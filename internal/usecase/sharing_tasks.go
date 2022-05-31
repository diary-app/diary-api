package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type SharingTask struct {
	DiaryID           uuid.UUID `json:"diaryID" db:"diary_id"`
	ReceiverUserID    uuid.UUID `json:"receiverUserID" db:"receiver_user_id"`
	EncryptedDiaryKey string    `json:"encryptedDiaryKey" db:"encrypted_diary_key"`
	SharedAt          time.Time `json:"sharedAt" db:"shared_at"`
}

type NewSharingTaskRequest struct {
	DiaryID           uuid.UUID `json:"diaryID" binding:"required"`
	ReceiverUserID    uuid.UUID `json:"receiverUserID" binding:"required"`
	EncryptedDiaryKey string    `json:"encryptedDiaryKey" binding:"required"`
	SharedAt          time.Time `json:"sharedAt" binding:"required"`
}

type SharingTasksListResponse struct {
	Items []SharingTask `json:"items"`
}

type SharingTasksUseCase interface {
	CreateSharingTask(ctx context.Context, request *NewSharingTaskRequest) error
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	DeleteSharingTask(ctx context.Context, diaryID uuid.UUID, receiverID uuid.UUID) error
}

type SharingTasksRepository interface {
	CreateSharingTask(ctx context.Context, sharingTask *SharingTask) error
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	DeleteSharingTask(ctx context.Context, diaryID uuid.UUID, receiverID uuid.UUID) error
}

// Errors

var (
	ErrUserAlreadyHasTaskForSameDiary = errors.New("user already has sharing task for the same diary")
	ErrUserAlreadyHasAccessToDiary    = errors.New("user already has access to the diary")
)
