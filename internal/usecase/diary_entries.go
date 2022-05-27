package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type DiaryEntriesUseCase interface {
	GetById(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, request CreateDiaryEntryRequest) (*DiaryEntry, error)
	UpdateContents(ctx context.Context, contentsChanges DiaryEntryContentsChangeList)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetEntries(ctx context.Context, request GetDiaryEntriesParams) ([]DiaryEntry, error)
}

type DiaryEntriesRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*DiaryEntry, error)
	Create(ctx context.Context, entry *DiaryEntry) (*DiaryEntry, error)
	UpdateContents(ctx context.Context, contentsChanges DiaryEntryContentsChangeList)
	Delete(ctx context.Context, id uuid.UUID) error
	GetEntries(ctx context.Context, request GetDiaryEntriesParams) ([]DiaryEntry, error)
}

// Domain models

type DiaryEntry struct {
	Id       uuid.UUID     `json:"id"`
	DiaryId  uuid.UUID     `json:"diaryId"`
	Name     string        `json:"name"`
	Date     time.Time     `json:"date"`
	Contents []interface{} `json:"contents"`
}

type DiaryEntryBlock struct {
	Id    uuid.UUID
	Value map[string]interface{}
}

// DTO

type CreateDiaryEntryRequest struct {
	DiaryId uuid.UUID `json:"diaryId,omitempty" binding:"required"`
	Name    string    `json:"name,omitempty" binding:"required"`
	Date    time.Time `json:"date" binding:"required"`
}

type GetDiaryEntriesParams struct {
	DiaryId *uuid.UUID `uri:"diaryId,omitempty"`
	Date    *time.Time `uri:"date,omitempty"`
}

type GetDiaryEntriesResponse struct {
	Items []DiaryEntry `json:"items"`
}

type UpdateDiaryEntryRequest struct {
	DiaryEntryId uuid.UUID
	Name         *string
	Date         *time.Time
}

type DiaryEntryContentsChangeList struct {
	Contents []DiaryEntryContentChangeRequest
}

type DiaryEntryContentChangeRequest struct {
	ChangeType ContentChangeType       `json:"changeType"`
	Id         *uuid.UUID              `json:"id"`
	DiaryId    *uuid.UUID              `json:"diaryId,omitempty"`
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
