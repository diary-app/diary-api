package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type SharingTasksUseCase interface {
	CreateSharingTask(ctx context.Context, req *NewSharingTaskRequest) (*SharingTask, error)
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	AcceptSharingTask(ctx context.Context, req *AcceptSharingTaskRequest) error
}

type SharingTasksRepository interface {
	CreateSharingTask(ctx context.Context, req *NewSharingTaskRequest) (*SharingTask, error)
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	AcceptSharingTask(ctx context.Context, req *AcceptSharingTaskRequest) error
}

// Models

type SharingTask struct {
	DiaryID           uuid.UUID `json:"diaryID" db:"diary_id"`
	ReceiverUserID    uuid.UUID `json:"receiverUserID" db:"receiver_user_id"`
	EncryptedDiaryKey string    `json:"encryptedDiaryKey" db:"encrypted_diary_key"`
	SharedAt          time.Time `json:"sharedAt" db:"shared_at"`
}

// DTO

type NewSharingTaskRequest struct {
	EntryID              uuid.UUID            `json:"entryId" binding:"required"`
	ReceiverUserID       uuid.UUID            `json:"receiverUserId" binding:"required"`
	MyEncryptedKey       string               `json:"myEncryptedKey" binding:"required"`
	ReceiverEncryptedKey string               `json:"receiverEncryptedKey" binding:"required"`
	Blocks               []DiaryEntryBlockDto `json:"blocks" binding:"required"`
}

type SharingTasksListResponse struct {
	Items []SharingTask `json:"items"`
}

type AcceptSharingTaskRequest struct {
	DiaryID           uuid.UUID `json:"diaryId" binding:"required"`
	EncryptedDiaryKey string    `json:"encryptedDiaryKey" binding:"required"`
}

// Errors

var (
	ErrUserAlreadyHasTaskForSameDiary = errors.New("user already has sharing task for the same diary")
	ErrUserAlreadyHasAccessToDiary    = errors.New("user already has access to the diary")
)
