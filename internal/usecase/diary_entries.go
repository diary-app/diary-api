package usecase

import (
	"context"
	"github.com/google/uuid"
	"io"
	"time"
)

type DiaryEntry struct {
	Id          uuid.UUID `json:"id"`
	DiaryId     uuid.UUID `json:"diaryId"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	ContentPath string    `json:"-"`
}

type CreateDiaryEntryRequest struct {
	DiaryId uuid.UUID `json:"diaryId,omitempty"`
	Name    string    `json:"name,omitempty"`
	Date    time.Time `json:"date"`
}

type GetDiaryEntriesRequest struct {
	DiaryId   *uuid.UUID
	EntryDate *time.Time
}

type UpdateDiaryEntryRequest struct {
	DiaryEntryId uuid.UUID
	Name         *string
}

type DiaryEntriesUseCase interface {
	Create(ctx context.Context, request CreateDiaryEntryRequest) (*DiaryEntry, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	DownloadContents(ctx context.Context, id uuid.UUID) (io.Reader, error)
	GetEntries(ctx context.Context, request GetDiaryEntriesRequest) ([]DiaryEntry, error)
	Update(ctx context.Context, request UpdateDiaryEntryRequest) (*DiaryEntry, error)
	Upload(ctx context.Context, id uuid.UUID, contentsStream io.Reader) error
}

type DiaryEntriesRepository interface {
	Create(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error)
	Read(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Update(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ReadMany(ctx context.Context, request GetDiaryEntriesRequest) ([]DiaryEntry, error)
}
