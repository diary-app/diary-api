package repository

import (
	"context"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgSharingTasksRepo struct {
	db *sqlx.DB
}

func (p *pgSharingTasksRepo) CreateSharingTask(ctx context.Context, sharingTask *usecase.SharingTask) error {
	const checkKeyQuery = `SELECT EXISTS(SELECT * FROM diary_keys WHERE diary_id = $1 AND user_id = $2)`
	var exists bool
	err := p.db.QueryRowxContext(ctx, checkKeyQuery, sharingTask.DiaryID, sharingTask.ReceiverUserID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return usecase.ErrUserAlreadyHasAccessToDiary
	}

	const query = `
INSERT INTO sharing_tasks (diary_id, receiver_user_id, encrypted_diary_key, shared_at) 
VALUES (:diary_id, :receiver_user_id, :encrypted_diary_key, :shared_at)`
	_, err = p.db.NamedExecContext(ctx, query, sharingTask)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == db.UniqueViolationErrorCode {
				return usecase.ErrUserAlreadyHasTaskForSameDiary
			}
		}
		return err
	}
	return nil
}

func (p *pgSharingTasksRepo) GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]usecase.SharingTask, error) {
	const query = `
SELECT diary_id, receiver_user_id, encrypted_diary_key, shared_at FROM sharing_tasks WHERE receiver_user_id = $1`
	tasksArr := make([]usecase.SharingTask, 0)
	if err := p.db.SelectContext(ctx, &tasksArr, query, userID); err != nil {
		return nil, err
	}

	return tasksArr, nil
}

func (p *pgSharingTasksRepo) DeleteSharingTask(ctx context.Context, diaryID uuid.UUID, receiverID uuid.UUID) error {
	const query = `
DELETE FROM sharing_tasks WHERE diary_id = $1 AND receiver_user_id = $2`

	if _, err := p.db.ExecContext(ctx, query, diaryID, receiverID); err != nil {
		return err
	}

	return nil
}

func New(db *sqlx.DB) usecase.SharingTasksRepository {
	return &pgSharingTasksRepo{db: db}
}
