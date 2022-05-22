package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

type Diary struct {
	Id      uuid.UUID  `json:"id"`
	Name    string     `json:"name"`
	OwnerId uuid.UUID  `json:"ownerId"`
	Keys    []DiaryKey `json:"keys"`
}

type DiaryKey struct {
	DiaryId      uuid.UUID `json:"diaryId"`
	UserId       uuid.UUID `json:"userId"`
	EncryptedKey string    `json:"encryptedKey"`
}

type CreateDiaryRequest struct {
	Name         string `json:"name" binding:"required"`
	EncryptedKey string `json:"encryptedKey" binding:"required"`
}

type DiaryUseCase interface {
	CreateDiary(ctx context.Context, userId uuid.UUID, req *CreateDiaryRequest) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}

type DiaryRepository interface {
	CreateDiary(ctx context.Context, diary *Diary) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}

// Errors

var (
	ErrDuplicateDiaryName = errors.New("user already has diary with this name")
)
