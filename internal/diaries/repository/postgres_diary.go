package repository

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
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
	Id           uuid.UUID `db:"id"`
	Name         string    `db:"name"`
	OwnerId      uuid.UUID `db:"owner_id"`
	CreatedAt    time.Time `db:"created_at"`
	EncryptedKey string    `db:"encrypted_key"`
}

func (p *postgresDiaryRepository) GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]usecase.Diary, error) {
	var diariesWithKeys []diaryWithKey
	if err := p.db.SelectContext(ctx, &diariesWithKeys,
		"SELECT d.id, d.name, d.owner_id, d.created_at, k.encrypted_key FROM diary_keys k JOIN diaries d ON d.id = k.diary_id WHERE user_id = $1",
		userId); err != nil {
		return nil, err
	}

	diaries := make([]usecase.Diary, len(diariesWithKeys))
	for i, d := range diariesWithKeys {
		diaries[i] = mapToDiary(d)
	}

	return diaries, nil
}

func mapToDiary(d diaryWithKey) usecase.Diary {
	keys := make([]usecase.DiaryKey, 1)
	keys[0] = usecase.DiaryKey{
		DiaryId:      d.Id,
		UserId:       d.OwnerId,
		EncryptedKey: d.EncryptedKey,
	}

	diary := usecase.Diary{
		Id:        d.Id,
		Name:      d.Name,
		OwnerId:   d.OwnerId,
		CreatedAt: d.CreatedAt,
		Keys:      keys,
	}

	return diary
}
