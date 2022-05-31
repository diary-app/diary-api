package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type DiaryEntriesUseCase interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, request CreateDiaryEntryRequest) (*DiaryEntry, error)
	UpdateContents(ctx context.Context, contentsChanges DiaryEntryContentsChangeList)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetEntries(ctx context.Context, request GetDiaryEntriesParams) ([]DiaryEntry, error)
}

type DiaryEntriesRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error)
	UpdateContents(ctx context.Context, contentsChanges DiaryEntryContentsChangeList)
	Delete(ctx context.Context, id uuid.UUID) error
	GetEntries(ctx context.Context, request GetDiaryEntriesParams) ([]DiaryEntry, error)
}

// Domain models

type DiaryEntry struct {
	ID       uuid.UUID     `json:"id"`
	DiaryID  uuid.UUID     `json:"diaryID"`
	Name     string        `json:"name"`
	Date     time.Time     `json:"date"`
	Contents []interface{} `json:"contents"`
}

type DiaryEntryBlock struct {
	ID    uuid.UUID
	Value map[string]interface{}
}

// DTO

type CreateDiaryEntryRequest struct {
	DiaryID uuid.UUID `json:"diaryID,omitempty" binding:"required"`
	Name    string    `json:"name,omitempty" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}

type GetDiaryEntriesParams struct {
	DiaryID *uuid.UUID `uri:"diaryID,omitempty"`
	Date    *time.Time `uri:"date,omitempty"`
}

type GetDiaryEntriesResponse struct {
	Items []DiaryEntry `json:"items"`
}

type UpdateDiaryEntryRequest struct {
	DiaryEntryID uuid.UUID
	Name         *string
	Date         *time.Time
}

type DiaryEntryContentsChangeList struct {
	Contents []DiaryEntryContentChangeRequest
}

type DiaryEntryContentChangeRequest struct {
	ChangeType ContentChangeType       `json:"changeType"`
	ID         *uuid.UUID              `json:"id"`
	DiaryID    *uuid.UUID              `json:"diaryID,omitempty"`
	Value      *map[string]interface{} `json:"value"`
}

type ContentChangeType string

const (
	Create ContentChangeType = "create"
	Update                   = "update"
	Delete                   = "delete"
)

// Errors

var (
	ErrNoAccessToDiary = errors.New("user does not have access to diary")
)
