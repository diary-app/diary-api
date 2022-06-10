package usecase

import (
	"context"
	"diary-api/internal/protocol/rest/common"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type DiaryEntriesUseCase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, req CreateDiaryEntryRequest) (*DiaryEntry, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateDiaryEntryRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEntries(ctx context.Context, req GetDiaryEntriesParams) ([]DiaryEntry, error)
}

type DiaryEntriesRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateDiaryEntryRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEntries(ctx context.Context, request GetDiaryEntriesParams) ([]DiaryEntry, error)
}

// Domain models

type DiaryEntry struct {
	ID      uuid.UUID         `json:"id" db:"id"`
	DiaryID uuid.UUID         `json:"diaryID" db:"diary_id"`
	Name    string            `json:"name" db:"name"`
	Date    time.Time         `json:"date" db:"date"`
	Value   string            `json:"value" db:"value"`
	Blocks  []DiaryEntryBlock `json:"blocks"`
}

type DiaryEntryBlock struct {
	ID    uuid.UUID
	Value string
}

// DTO

type CreateDiaryEntryRequest struct {
	DiaryID uuid.UUID       `json:"diaryId" binding:"required"`
	Name    string          `json:"name" binding:"required"`
	Date    common.DateOnly `json:"date" binding:"required"`
	Value   string          `json:"value" binding:"required"`
}

type UpdateDiaryEntryRequest struct {
	DiaryId        *uuid.UUID           `json:"diaryId"`
	Name           *string              `json:"name"`
	Date           *common.DateOnly     `json:"date"`
	Value          *string              `json:"value"`
	BlocksToUpsert []DiaryEntryBlockDto `json:"blocksToUpsert"`
	BlocksToDelete []uuid.UUID          `json:"blocksToDelete"`
}

type DiaryEntryBlockDto struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Value string
}

type GetDiaryEntriesParams struct {
	DiaryID *uuid.UUID       `uri:"diaryID,omitempty"`
	Date    *common.DateOnly `uri:"date,omitempty"`
}

type ShortDiaryEntryResponse struct {
	ID      uuid.UUID       `json:"Id"`
	DiaryID uuid.UUID       `json:"diaryId"`
	Name    string          `json:"name"`
	Date    common.DateOnly `json:"date"`
}

type DiaryEntryResponse struct {
	ID      uuid.UUID                 `json:"Id"`
	DiaryID uuid.UUID                 `json:"diaryId"`
	Name    string                    `json:"name"`
	Date    common.DateOnly           `json:"date"`
	Blocks  []DiaryEntryBlockResponse `json:"blocks"`
}

type DiaryEntryBlockResponse struct {
	ID    uuid.UUID `json:"id"`
	Value string    `json:"value"`
}

// Errors

type NoAccessToDiaryEntryError struct {
	EntryID uuid.UUID
}

func (e *NoAccessToDiaryEntryError) Error() string {
	return fmt.Sprintf("no access to entry %v", e.EntryID)
}

type NoAccessToDiaryError struct {
	DiaryID uuid.UUID
}

func (e *NoAccessToDiaryError) Error() string {
	return fmt.Sprintf("no access to diary %v", e.DiaryID)
}

type AlienEntryBlocksError struct {
	DiaryEntryId   uuid.UUID
	AlienBlocksIds []uuid.UUID
}

func (e *AlienEntryBlocksError) Error() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("entry %s does not contains these blocks: ", e.DiaryEntryId))
	for _, alienBlocksId := range e.AlienBlocksIds {
		b.WriteString(fmt.Sprintf(", %s", alienBlocksId))
	}
	return b.String()
}
