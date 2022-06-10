package repository

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgSharingTasksRepo struct {
	db    db.Db
	clock clock.Clock
}

func (p *pgSharingTasksRepo) CreateSharingTask(
	ctx context.Context,
	req *usecase.NewSharingTaskRequest,
) (*usecase.SharingTask, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	userId := auth.MustGetUserID(ctx)
	if err = db.CheckUserAccessToDiary(ctx, tx, req.EntryID, userId); err != nil {
		return nil, err
	}

	err = db.CheckUserAccessToEntry(ctx, tx, req.EntryID, req.ReceiverUserID)
	if err == nil {
		return nil, usecase.ErrUserAlreadyHasAccessToDiary
	}
	if _, ok := err.(*usecase.NoAccessToDiaryEntryError); !ok {
		return nil, err
	}

	newDiaryID, err := moveEntryToNewDiary(ctx, tx, req.ReceiverUserID, req.EntryID)
	if err != nil {
		return nil, err
	}

	task, err := p.createSharingTask(ctx, tx, req, newDiaryID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func moveEntryToNewDiary(
	ctx context.Context,
	tx db.TxOrDb,
	receiverID uuid.UUID,
	entryID uuid.UUID,
) (uuid.UUID, error) {

	userID := auth.MustGetUserID(ctx)
	owner, err := getUser(ctx, tx, userID)
	if err != nil {
		return uuid.UUID{}, nil
	}
	receiver, err := getUser(ctx, tx, receiverID)
	if err != nil {
		return uuid.UUID{}, nil
	}

	//language=postgresql
	const createDiaryQuery = `
INSERT INTO diaries (id, name, owner_id) VALUES (:id, :name, :owner_id)`
	newDiaryID, _ := uuid.NewUUID()
	args := map[string]interface{}{
		"id":       newDiaryID,
		"name":     fmt.Sprintf("%s -> %s", owner.Username, receiver.Username),
		"owner_id": owner.ID,
	}
	if _, err = tx.NamedExecContext(ctx, createDiaryQuery, args); err != nil {
		return uuid.UUID{}, nil
	}

	//language=postgresql
	const moveEntryQuery = `UPDATE diary_entries SET diary_id = $2 WHERE id = $1`
	if _, err = tx.ExecContext(ctx, moveEntryQuery, entryID, newDiaryID); err != nil {
		return uuid.UUID{}, err
	}

	return newDiaryID, nil
}

func getUser(ctx context.Context, tx db.TxOrDb, id uuid.UUID) (*usecase.FullUser, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	user := &usecase.FullUser{}
	if err := tx.GetContext(ctx, user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *pgSharingTasksRepo) createSharingTask(
	ctx context.Context,
	tx db.TxOrDb,
	req *usecase.NewSharingTaskRequest,
	newDiaryID uuid.UUID,
) (*usecase.SharingTask, error) {

	const query = `
INSERT INTO sharing_tasks (diary_id, receiver_user_id, encrypted_diary_key, shared_at) 
VALUES (:diary_id, :receiver_user_id, :encrypted_diary_key, :shared_at)`
	task := &usecase.SharingTask{
		DiaryID:           newDiaryID,
		ReceiverUserID:    req.ReceiverUserID,
		EncryptedDiaryKey: req.ReceiverEncryptedKey,
		SharedAt:          r.clock.Now(),
	}
	if _, err := tx.NamedExecContext(ctx, query, task); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == db.UniqueViolationErrorCode {
				return nil, usecase.ErrUserAlreadyHasTaskForSameDiary
			}
		}
		return nil, err
	}
	return task, nil
}

func (p *pgSharingTasksRepo) GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]usecase.SharingTask, error) {
	//language=postgresql
	const query = `
SELECT diary_id, receiver_user_id, encrypted_diary_key, shared_at FROM sharing_tasks WHERE receiver_user_id = $1`
	tasksArr := make([]usecase.SharingTask, 0)
	if err := p.db.SelectContext(ctx, &tasksArr, query, userID); err != nil {
		return nil, err
	}

	return tasksArr, nil
}

func (p *pgSharingTasksRepo) AcceptSharingTask(ctx context.Context, req *usecase.AcceptSharingTaskRequest) error {
	//	tx, err := p.db.Beginx()
	//	if err != nil {
	//		return err
	//	}
	//
	//	userID := auth.MustGetUserID(ctx)
	//	//language=postgresql
	//	const insertKeyQuery = `
	//INSERT INTO diary_keys (diary_id, user_id, encrypted_key) VALUES (:diary_id, :user_id, :encrypted_key)`
	//
	//	//language=postgresql
	//	const query = `
	//DELETE FROM sharing_tasks WHERE diary_id = $1 AND receiver_user_id = $2`
	//
	//	if _, err := p.db.ExecContext(ctx, query, diaryID, receiverID); err != nil {
	//		return err
	//	}

	return nil
}

func New(db *sqlx.DB) usecase.SharingTasksRepository {
	return &pgSharingTasksRepo{db: db}
}
