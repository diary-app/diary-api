package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type Diary struct {
	ID      uuid.UUID  `json:"id"`
	Name    string     `json:"name"`
	OwnerID uuid.UUID  `json:"ownerID"`
	Keys    []DiaryKey `json:"keys"`
}

type DiaryKey struct {
	DiaryID      uuid.UUID `json:"diaryID"`
	UserID       uuid.UUID `json:"userID"`
	EncryptedKey string    `json:"encryptedKey"`
}

type CreateDiaryRequest struct {
	Name         string `json:"name" binding:"required"`
	EncryptedKey string `json:"encryptedKey" binding:"required"`
}

type DiaryUseCase interface {
	GetDiariesByUser(ctx context.Context) ([]Diary, error)
}

type DiaryRepository interface {
	CreateDiary(ctx context.Context, diary *Diary) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userID uuid.UUID) ([]Diary, error)
}

// Errors

var (
	ErrDuplicateDiaryName = errors.New("user already has diary with this name")
)
