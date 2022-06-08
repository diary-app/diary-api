package usecase

import (
	"context"
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
	ID      uuid.UUID         `json:"id"`
	DiaryID uuid.UUID         `json:"diaryID"`
	Name    string            `json:"name"`
	Date    time.Time         `json:"date"`
	Blocks  []DiaryEntryBlock `json:"blocks,omitempty"`
}

type DiaryEntryBlock struct {
	ID    uuid.UUID
	Value string
}

// DTO

type CreateDiaryEntryRequest struct {
	DiaryID uuid.UUID `json:"diaryId" binding:"required"`
	Name    string    `json:"name" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
	Value   string    `json:"value" binding:"required"`
}

type UpdateDiaryEntryRequest struct {
	DiaryId        *uuid.UUID           `json:"diaryId" binding:"optional"`
	Name           *string              `json:"name" binding:"optional"`
	Date           *time.Time           `json:"date" binding:"optional"`
	Value          *string              `json:"value" binding:"optional"`
	BlocksToUpsert []DiaryEntryBlockDto `json:"blocksToUpsert" binding:"optional"`
	BlocksToDelete []uuid.UUID          `json:"blocksToDelete" binding:"optional"`
}

type DiaryEntryBlockDto struct {
	ID    uuid.UUID `json:"id" binding:"required"`
	Value string
}

type GetDiaryEntriesParams struct {
	DiaryID *uuid.UUID `uri:"diaryID,omitempty" binding:"optional"`
	Date    *time.Time `uri:"date,omitempty" binding:"optional"`
}

type ShortDiaryEntryResponse struct {
	ID      uuid.UUID `json:"Id"`
	DiaryID uuid.UUID `json:"diaryId"`
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
}

type DiaryEntryResponse struct {
	ID      uuid.UUID                 `json:"Id"`
	DiaryID uuid.UUID                 `json:"diaryId"`
	Name    string                    `json:"name"`
	Date    time.Time                 `json:"date"`
	Blocks  []DiaryEntryBlockResponse `json:"blocks"`
}

type DiaryEntryBlockResponse struct {
	ID    uuid.UUID `json:"id"`
	Value string    `json:"value"`
}

// Errors

type NoAccessToDiaryEntryError struct {
	DiaryId uuid.UUID
}

func (e *NoAccessToDiaryEntryError) Error() string {
	return fmt.Sprintf("user does not have access to diary %v", e.DiaryId)
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
