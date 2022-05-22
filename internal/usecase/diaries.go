package usecase

import (
	"context"
	"github.com/google/uuid"
)

type Diary struct {
	Id      uuid.UUID  `json:"id,omitempty"`
	Name    string     `json:"name,omitempty"`
	OwnerId uuid.UUID  `json:"ownerId,omitempty"`
	Keys    []DiaryKey `json:"keys,omitempty"`
}

type DiaryKey struct {
	DiaryId      uuid.UUID `json:"diaryId,omitempty"`
	UserId       uuid.UUID `json:"userId,omitempty"`
	EncryptedKey string    `json:"encryptedKey,omitempty"`
}

type CreateDiaryRequest struct {
	Name         string
	EncryptedKey string
}

type DiaryUseCase interface {
	CreateDiary(ctx context.Context, userId uuid.UUID, req *CreateDiaryRequest) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}

type DiaryRepository interface {
	CreateDiary(ctx context.Context, diary *Diary) (*Diary, error)
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}
