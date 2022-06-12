package repository

import (
	"context"
	"database/sql"
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

func (r *pgSharingTasksRepo) CreateSharingTask(
	ctx context.Context,
	req *usecase.CreateSharingTaskRequest,
) (*usecase.SharingTask, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	userId := auth.MustGetUserID(ctx)
	if err = db.CheckUserAccessToEntry(ctx, tx, req.EntryID, userId); err != nil {
		return nil, db.ShouldRollback(tx, err)
	}

	err = db.CheckUserAccessToEntry(ctx, tx, req.EntryID, req.ReceiverUserID)
	if err == nil {
		return nil, db.ShouldRollback(tx, usecase.ErrUserAlreadyHasAccessToDiary)
	}
	if _, ok := err.(*usecase.NoAccessToDiaryEntryError); !ok {
		return nil, db.ShouldRollback(tx, err)
	}

	if err = checkSharingTask(ctx, tx, req.EntryID, req.ReceiverUserID); err != nil {
		return nil, db.ShouldRollback(tx, err)
	}

	newDiaryID, err := moveEntryToNewDiary(ctx, tx, req)
	if err != nil {
		return nil, db.ShouldRollback(tx, err)
	}

	err = updateBlocks(ctx, tx, req.EntryID, req.Blocks)
	if err != nil {
		return nil, db.ShouldRollback(tx, err)
	}

	task, err := r.createSharingTask(ctx, tx, req, newDiaryID)
	if err != nil {
		return nil, db.ShouldRollback(tx, err)
	}

	if err = db.ShouldCommitOrRollback(tx); err != nil {
		return nil, err
	}
	return task, nil
}

func checkSharingTask(ctx context.Context, tx db.TxOrDb, entryID uuid.UUID, receiverUserID uuid.UUID) error {
	const checkSharingTaskQuery = `
SELECT EXISTS(SELECT 1 
              FROM diary_entries de 
                  JOIN diaries d ON de.diary_id = d.id 
                  JOIN sharing_tasks st ON d.id = st.diary_id 
              WHERE de.id = $1 AND st.receiver_user_id = $2)`

	var taskExists bool
	if err := tx.QueryRowxContext(ctx, checkSharingTaskQuery, entryID, receiverUserID).Scan(&taskExists); err != nil {
		return err
	}
	if taskExists {
		return usecase.ErrUserAlreadyHasTaskForSameDiary
	}
	return nil
}

func updateBlocks(ctx context.Context, tx db.TxOrDb, entryID uuid.UUID, blocks []usecase.DiaryEntryBlockDto) error {
	const getAllBlocksIdsQuery = `SELECT id FROM diary_entry_blocks WHERE diary_entry_id = $1`
	var existingIds []uuid.UUID
	if err := tx.SelectContext(ctx, &existingIds, getAllBlocksIdsQuery, entryID); err != nil {
		return err
	}

	seenIds := map[uuid.UUID]bool{}
	for _, ids := range existingIds {
		seenIds[ids] = false
	}
	var alienIds []uuid.UUID
	var duplicatedIds []uuid.UUID
	for _, b := range blocks {
		if _, ok := seenIds[b.ID]; !ok {
			alienIds = append(alienIds, b.ID)
		} else if !seenIds[b.ID] {
			seenIds[b.ID] = true
		} else {
			duplicatedIds = append(duplicatedIds, b.ID)
		}
	}
	var missingIds []uuid.UUID
	for id, seen := range seenIds {
		if !seen {
			missingIds = append(missingIds, id)
		}
	}

	if len(alienIds) > 0 || len(duplicatedIds) > 0 || len(missingIds) > 0 {
		return &usecase.BadUpdatedBlocksError{
			AlienBlocks: alienIds, DuplicatedBlocks: duplicatedIds, MissingBlocks: missingIds}
	}

	entriesToUpdate := make([]diaryEntryBlock, len(blocks))
	for i, b := range blocks {
		entriesToUpdate[i] = diaryEntryBlock{ID: b.ID, Value: b.Value}
	}
	const updateQuery = `UPDATE diary_entry_blocks SET value = :value WHERE id = :id`
	for _, entry := range entriesToUpdate {
		if _, err := tx.NamedExecContext(ctx, updateQuery, entry); err != nil {
			return err
		}
	}
	return nil
}

type diaryEntryBlock struct {
	ID    uuid.UUID `db:"id"`
	Value []byte    `db:"value"`
}

func moveEntryToNewDiary(
	ctx context.Context,
	tx db.TxOrDb,
	req *usecase.CreateSharingTaskRequest,
) (uuid.UUID, error) {

	userID := auth.MustGetUserID(ctx)
	owner, err := getUser(ctx, tx, userID)
	if err != nil {
		return uuid.UUID{}, nil
	}
	receiver, err := getUser(ctx, tx, req.ReceiverUserID)
	if err == sql.ErrNoRows {
		return uuid.UUID{}, usecase.ErrReceiverUserNotFound
	}
	if err != nil {
		return uuid.UUID{}, err
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
		return uuid.UUID{}, err
	}

	//language=postgresql
	const createDiaryKeyQuery = `INSERT INTO diary_keys (diary_id, user_id, encrypted_key) VALUES ($1, $2, $3)`
	if _, err = tx.ExecContext(ctx, createDiaryKeyQuery, newDiaryID, userID, req.MyEncryptedKey, req.Value); err != nil {
		return uuid.UUID{}, err
	}

	//language=postgresql
	const moveEntryQuery = `UPDATE diary_entries SET diary_id = $2, value = $3 WHERE id = $1`
	if _, err = tx.ExecContext(ctx, moveEntryQuery, req.EntryID, newDiaryID); err != nil {
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
	req *usecase.CreateSharingTaskRequest,
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

func (r *pgSharingTasksRepo) GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]usecase.SharingTask, error) {
	//language=postgresql
	const query = `
SELECT st.diary_id, st.receiver_user_id, st.encrypted_diary_key, st.shared_at, u.username
FROM sharing_tasks st
JOIN diaries d ON st.diary_id = d.id
JOIN users u ON d.owner_id = u.id
WHERE receiver_user_id = $1`
	tasksArr := make([]usecase.SharingTask, 0)
	if err := r.db.SelectContext(ctx, &tasksArr, query, userID); err != nil {
		return nil, err
	}

	return tasksArr, nil
}

func (r *pgSharingTasksRepo) AcceptSharingTask(ctx context.Context, req *usecase.AcceptSharingTaskRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	userID := auth.MustGetUserID(ctx)
	//language=postgresql
	const checkTaskQuery = `SELECT EXISTS(
    	SELECT * FROM sharing_tasks 
		WHERE diary_id = $1 AND receiver_user_id = $2)`
	var taskExists bool
	if err = tx.QueryRowxContext(ctx, checkTaskQuery, req.DiaryID, userID).Scan(&taskExists); err != nil {
		return db.ShouldRollback(tx, err)
	}
	if !taskExists {
		return db.ShouldRollback(tx, usecase.ErrCommonNotFound)
	}

	//language=postgresql
	const insertKeyQuery = `
	INSERT INTO diary_keys (diary_id, user_id, encrypted_key) VALUES ($1, $2, $3)`
	if _, err = tx.ExecContext(ctx, insertKeyQuery, req.DiaryID, userID, req.EncryptedDiaryKey); err != nil {
		return db.ShouldRollback(tx, err)
	}

	//language=postgresql
	const query = `
	DELETE FROM sharing_tasks WHERE diary_id = $1 AND receiver_user_id = $2`
	if _, err = tx.ExecContext(ctx, query, req.DiaryID, userID); err != nil {
		return db.ShouldRollback(tx, err)
	}

	if err = db.ShouldCommitOrRollback(tx); err != nil {
		return err
	}

	return nil
}

func New(db *sqlx.DB, clock clock.Clock) usecase.SharingTasksRepository {
	return &pgSharingTasksRepo{db: db, clock: clock}
}
