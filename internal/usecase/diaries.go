package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type DiaryUseCase interface {
	GetDiariesByUser(ctx context.Context) ([]Diary, error)
}

type DiaryRepository interface {
	CreateDiary(ctx context.Context, diary *Diary) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userID uuid.UUID) ([]Diary, error)
}

// Models

type Diary struct {
	ID      uuid.UUID  `json:"id" db:"id"`
	Name    string     `json:"name" db:"name"`
	OwnerID uuid.UUID  `json:"ownerID" db:"owner_id"`
	Keys    []DiaryKey `json:"keys"`
}

type DiaryKey struct {
	DiaryID      uuid.UUID `json:"diaryID" db:"diary_id"`
	UserID       uuid.UUID `json:"userID" db:"user_id"`
	EncryptedKey string    `json:"encryptedKey" db:"encrypted_key"`
}

// DTO

type CreateDiaryRequest struct {
	Name         string `json:"name" binding:"required"`
	EncryptedKey string `json:"encryptedKey" binding:"required"`
}

type DiaryResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	OwnerID      uuid.UUID `json:"ownerID"`
	EncryptedKey string    `json:"encryptedKey"`
}

type DiaryListResponse struct {
	Items []DiaryResponse `json:"items"`
}

// Errors

var (
	ErrDuplicateDiaryName = errors.New("user already has diary with this name")
)
