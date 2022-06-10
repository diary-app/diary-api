package repository

import (
	"context"
	"diary-api/internal/db"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type postgresDiaryRepository struct {
	db *sqlx.DB
}

func NewPostgresDiaryRepository(db *sqlx.DB) usecase.DiaryRepository {
	return &postgresDiaryRepository{
		db: db,
	}
}

type diaryWithKey struct {
	ID           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	OwnerID      uuid.UUID `db:"owner_id"`
	EncryptedKey string    `db:"encrypted_key"`
}

type newDiaryKey struct {
	DiaryID      uuid.UUID `db:"diary_id"`
	UserID       uuid.UUID `db:"user_id"`
	EncryptedKey string    `db:"encrypted_key"`
}

type newDiary struct {
	OwnerID uuid.UUID `db:"owner_id"`
	Name    string    `db:"name"`
}

func (p *postgresDiaryRepository) CreateDiary(ctx context.Context, diary *usecase.Diary) (*usecase.Diary, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	diary, err = insertDiary(ctx, diary, tx)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == db.UniqueViolationErrorCode {
			return nil, usecase.ErrDuplicateDiaryName
		}
		return nil, err
	}

	if diary.Keys != nil && len(diary.Keys) > 0 {
		for i := range diary.Keys {
			diary.Keys[i].DiaryID = diary.ID
		}

		if err = insertDiaryKeys(ctx, diary.Keys, tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return nil, multierror.Append(err, rbErr)
			}
			return nil, err
		}
	}

	if err = db.ShouldCommitOrRollback(tx); err != nil {
		return nil, err
	}
	return diary, nil
}

func (p *postgresDiaryRepository) GetDiariesByUser(ctx context.Context, userID uuid.UUID) ([]usecase.Diary, error) {
	const query = `
SELECT d.id, d.name, d.owner_id, k.encrypted_key
FROM diary_keys k
    JOIN diaries d ON d.id = k.diary_id
WHERE k.user_id = $1
`
	var diariesWithKeys []diaryWithKey
	if err := p.db.SelectContext(ctx, &diariesWithKeys, query, userID); err != nil {
		return nil, err
	}

	diaries := make([]usecase.Diary, len(diariesWithKeys))
	for i, d := range diariesWithKeys {
		keys := make([]usecase.DiaryKey, 1)
		keys[0] = usecase.DiaryKey{
			DiaryID:      d.ID,
			UserID:       d.OwnerID,
			EncryptedKey: d.EncryptedKey,
		}

		diary := usecase.Diary{
			ID:      d.ID,
			Name:    d.Name,
			OwnerID: d.OwnerID,
			Keys:    keys,
		}
		diaries[i] = diary
	}

	return diaries, nil
}

func insertDiary(ctx context.Context, diary *usecase.Diary, tx *sqlx.Tx) (*usecase.Diary, error) {
	newD := newDiary{
		OwnerID: diary.OwnerID,
		Name:    diary.Name,
	}
	const diaryQuery = `INSERT INTO diaries(name, owner_id) VALUES(:name,:owner_id) RETURNING id`
	query, args, err := tx.BindNamed(diaryQuery, newD)
	if err != nil {
		return nil, err
	}
	var diaryID uuid.UUID
	err = tx.GetContext(ctx, &diaryID, query, args...)
	if err != nil {
		return nil, err
	}
	diary.ID = diaryID
	return diary, nil
}

func insertDiaryKeys(ctx context.Context, keys []usecase.DiaryKey, tx *sqlx.Tx) error {
	const diaryKeysQuery = `
INSERT INTO diary_keys (diary_id, user_id, encrypted_key) VALUES (:diary_id, :user_id, :encrypted_key)`
	newDiaryKeys := make([]newDiaryKey, len(keys))
	for i, key := range keys {
		newDiaryKeys[i] = newDiaryKey{
			DiaryID:      key.DiaryID,
			UserID:       key.UserID,
			EncryptedKey: key.EncryptedKey,
		}
	}

	query, args, err := tx.BindNamed(diaryKeysQuery, newDiaryKeys)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
