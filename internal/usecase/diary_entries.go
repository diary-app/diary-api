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
	Value   []byte            `json:"value" db:"value"`
	Blocks  []DiaryEntryBlock `json:"blocks"`
}

type DiaryEntryBlock struct {
	ID    uuid.UUID `json:"id"`
	Value []byte    `json:"value"`
}

// DTO

type CreateDiaryEntryRequest struct {
	DiaryID uuid.UUID       `json:"diaryId" binding:"required"`
	Name    string          `json:"name" binding:"required"`
	Date    common.DateOnly `json:"date" binding:"required"`
	Value   []byte          `json:"value"`
}

type UpdateDiaryEntryRequest struct {
	DiaryId        *uuid.UUID           `json:"diaryId"`
	Name           *string              `json:"name"`
	Date           *common.DateOnly     `json:"date"`
	Value          []byte               `json:"value"`
	BlocksToUpsert []DiaryEntryBlockDto `json:"blocksToUpsert"`
	BlocksToDelete []uuid.UUID          `json:"blocksToDelete"`
}

type DiaryEntryBlockDto struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Value []byte    `json:"value" binding:"required"`
}

type GetDiaryEntriesParamsDto struct {
	DiaryIDStr *string `form:"diaryId" binding:"omitempty"`
	DateStr    *string `form:"date" binding:"omitempty" time_format:"2006-01-02"`
}

type GetDiaryEntriesParams struct {
	DiaryID *uuid.UUID
	Date    *common.DateOnly
}

type ShortDiaryEntryResponse struct {
	ID      uuid.UUID       `json:"id"`
	DiaryID uuid.UUID       `json:"diaryId"`
	Name    string          `json:"name"`
	Date    common.DateOnly `json:"date"`
}

type DiaryEntryResponse struct {
	ID      uuid.UUID                 `json:"id"`
	DiaryID uuid.UUID                 `json:"diaryId"`
	Name    string                    `json:"name"`
	Date    common.DateOnly           `json:"date"`
	Value   []byte                    `json:"value"`
	Blocks  []DiaryEntryBlockResponse `json:"blocks"`
}

type DiaryEntryBlockResponse struct {
	ID    uuid.UUID `json:"id"`
	Value []byte    `json:"value"`
}

// Errors

type NoReadAccessToDiaryEntryError struct {
	EntryID uuid.UUID
}

func (e *NoReadAccessToDiaryEntryError) Error() string {
	return fmt.Sprintf("no read access to entry %v", e.EntryID)
}

type NoReadAccessToDiaryError struct {
	DiaryID uuid.UUID
}

func (e *NoReadAccessToDiaryError) Error() string {
	return fmt.Sprintf("no read access to diary %v", e.DiaryID)
}

type NoWriteAccessToDiaryError struct {
	DiaryID uuid.UUID
}

func (e *NoWriteAccessToDiaryError) Error() string {
	return fmt.Sprintf("no write access to diary %v", e.DiaryID)
}

type NoWriteAccessToDiaryEntryError struct {
	EntryID uuid.UUID
}

func (e *NoWriteAccessToDiaryEntryError) Error() string {
	return fmt.Sprintf("no write access to entry %v", e.EntryID)
}

type AlienEntryBlocksError struct {
	AlienBlocksIds []uuid.UUID
}

func (e *AlienEntryBlocksError) Error() string {
	var b strings.Builder
	b.WriteString("entry does not contains these blocks: ")
	for _, alienBlocksId := range e.AlienBlocksIds {
		b.WriteString(fmt.Sprintf("%s, ", alienBlocksId))
	}
	return strings.Trim(b.String(), ", ")
}
