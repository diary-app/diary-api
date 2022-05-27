package diary_entries

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"time"
)

type diaryEntriesUseCase struct {
	repo usecase.DiaryEntriesRepository
}

func (d *diaryEntriesUseCase) GetEntries(ctx context.Context, request usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error) {
	return d.repo.GetEntries(ctx, request)
}

func (d *diaryEntriesUseCase) GetById(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error) {
	return d.repo.GetById(ctx, id)
}

func (d *diaryEntriesUseCase) UpdateContents(ctx context.Context, contentsChanges usecase.DiaryEntryContentsChangeList) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) Create(ctx context.Context, r usecase.CreateDiaryEntryRequest) (*usecase.DiaryEntry, error) {
	id := uuid.New()
	date := r.Date
	entry := &usecase.DiaryEntry{
		Id:      id,
		DiaryId: r.DiaryId,
		Name:    r.Name,
		Date:    time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
	}
	entry, err := d.repo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func NewUseCase(repo usecase.DiaryEntriesRepository) usecase.DiaryEntriesUseCase {
	return &diaryEntriesUseCase{repo: repo}
}
