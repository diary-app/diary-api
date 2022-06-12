package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type SharingTasksUseCase interface {
	CreateSharingTask(ctx context.Context, req *CreateSharingTaskRequest) (*SharingTask, error)
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	AcceptSharingTask(ctx context.Context, req *AcceptSharingTaskRequest) error
}

type SharingTasksRepository interface {
	CreateSharingTask(ctx context.Context, req *CreateSharingTaskRequest) (*SharingTask, error)
	GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]SharingTask, error)
	AcceptSharingTask(ctx context.Context, req *AcceptSharingTaskRequest) error
}

// Models

type SharingTask struct {
	DiaryID           uuid.UUID `json:"diaryId" db:"diary_id"`
	ReceiverUserID    uuid.UUID `json:"receiverUserId" db:"receiver_user_id"`
	EncryptedDiaryKey []byte    `json:"encryptedDiaryKey" db:"encrypted_diary_key"`
	Username          string    `json:"username" db:"username"`
	SharedAt          time.Time `json:"sharedAt" db:"shared_at"`
}

// DTO

type CreateSharingTaskRequest struct {
	EntryID              uuid.UUID            `json:"entryId" binding:"required"`
	ReceiverUserID       uuid.UUID            `json:"receiverUserId" binding:"required"`
	MyEncryptedKey       []byte               `json:"myEncryptedKey" binding:"required"`
	ReceiverEncryptedKey []byte               `json:"receiverEncryptedKey" binding:"required"`
	Value                []byte               `json:"value" binding:"required"`
	Blocks               []DiaryEntryBlockDto `json:"blocks" binding:"required,dive"`
}

type CreateSharingTaskResponse struct {
	DiaryID uuid.UUID `json:"diaryId"`
}

type SharingTasksListResponse struct {
	Items []SharingTask `json:"items"`
}

type AcceptSharingTaskRequest struct {
	DiaryID           uuid.UUID `json:"diaryId" binding:"required"`
	EncryptedDiaryKey []byte    `json:"encryptedDiaryKey" binding:"required"`
}

// Errors

var (
	ErrUserAlreadyHasTaskForSameDiary = errors.New("user already has sharing task for the same diary")
	ErrUserAlreadyHasAccessToDiary    = errors.New("user already has access to the diary")
	ErrReceiverUserNotFound           = errors.New("receiver user not found")
)

type BadUpdatedBlocksError struct {
	AlienBlocks      []uuid.UUID
	DuplicatedBlocks []uuid.UUID
	MissingBlocks    []uuid.UUID
}

func (e *BadUpdatedBlocksError) Error() string {
	var sb strings.Builder
	if len(e.AlienBlocks) > 0 {
		sb.WriteString(fmt.Sprintf("blocks not contained in the entry: %s; ", idsToString(e.AlienBlocks)))
	}
	if len(e.DuplicatedBlocks) > 0 {
		sb.WriteString(fmt.Sprintf("blocks duplicated in the request: %s;", idsToString(e.DuplicatedBlocks)))
	}
	if len(e.MissingBlocks) > 0 {
		sb.WriteString(fmt.Sprintf("blocks missing in the requrest: %s", idsToString(e.MissingBlocks)))
	}
	return sb.String()
}

func idsToString(ids []uuid.UUID) string {
	var b strings.Builder
	for _, id := range ids {
		b.WriteString(fmt.Sprintf("%s, ", id))
	}
	return strings.Trim(b.String(), ", ")
}
