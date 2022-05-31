package usecase

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"time"
)

type UseCase struct {
	repo storage
}

func (d *UseCase) GetEntries(ctx context.Context, request usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error) {
	return d.repo.GetEntries(ctx, request)
}

func (d *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	return d.repo.GetByID(ctx, id)
}

func (d *UseCase) UpdateContents(ctx context.Context, contentsChanges usecase.DiaryEntryContentsChangeList) {
	//TODO implement me
	panic("implement me")
}

func (d *UseCase) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *UseCase) Create(ctx context.Context, r usecase.CreateDiaryEntryRequest) (*usecase.DiaryEntry, error) {
	id := uuid.New()
	date := r.Date
	entry := &usecase.DiaryEntry{
		ID:      id,
		DiaryID: r.DiaryID,
		Name:    r.Name,
		Date:    time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
	}
	entry, err := d.repo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func New(repo storage) *UseCase {
	return &UseCase{repo: repo}
}
