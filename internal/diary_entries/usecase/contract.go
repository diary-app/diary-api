//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package usecase

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type storage interface {
	GetByID(ctx context.Context, id uuid.UUID) (*usecase.DiaryEntry, error)
	Create(ctx context.Context, entry *usecase.DiaryEntry) (*usecase.DiaryEntry, error)
	Update(ctx context.Context, id uuid.UUID, r *usecase.UpdateDiaryEntryRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEntries(ctx context.Context, request usecase.GetDiaryEntriesParams) ([]usecase.DiaryEntry, error)
}
