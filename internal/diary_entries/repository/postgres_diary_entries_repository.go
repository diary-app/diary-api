package repository

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type pgRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) usecase.DiaryEntriesRepository {
	return &pgRepo{
		db: db,
	}
}

func (p *pgRepo) Create(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error) {
	_, err := p.db.NamedExecContext(ctx,
		"INSERT INTO diary_entries (id, diary_id, name, date, content_name) VALUES (:id, :diary_id, :name, :date, :content_path)", entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (p *pgRepo) Read(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepo) Update(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepo) ReadMany(ctx context.Context, request usecase.GetDiaryEntriesRequest) ([]usecase.DiaryEntry, error) {
	//TODO implement me
	panic("implement me")
}
