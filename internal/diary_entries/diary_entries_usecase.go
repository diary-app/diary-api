package diary_entries

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"io"
)

type diaryEntriesUseCase struct {
	repo usecase.DiaryEntriesRepository
}

func (d *diaryEntriesUseCase) Create(ctx context.Context, request usecase.CreateDiaryEntryRequest) (*usecase.DiaryEntry, error) {
	id := uuid.New()
	entry := &usecase.DiaryEntry{
		Id:          id,
		DiaryId:     request.DiaryId,
		Name:        request.Name,
		Date:        request.Date,
		ContentPath: "diary-entries/" + id.String(),
	}
	entry, err := d.repo.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (d *diaryEntriesUseCase) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) DownloadContents(ctx context.Context, id uuid.UUID) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) GetEntries(ctx context.Context, request usecase.GetDiaryEntriesRequest) ([]usecase.DiaryEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) Update(ctx context.Context, request usecase.UpdateDiaryEntryRequest) (*usecase.DiaryEntry, error) {
	//TODO implement me
	panic("implement me")
}

func (d *diaryEntriesUseCase) Upload(ctx context.Context, id uuid.UUID, contentsStream io.Reader) error {
	//TODO implement me
	panic("implement me")
}

func NewUseCase(repo usecase.DiaryEntriesRepository) usecase.DiaryEntriesUseCase {
	return &diaryEntriesUseCase{repo: repo}
}
