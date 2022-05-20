package usecase

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Diary struct {
	Id        uuid.UUID  `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	OwnerId   uuid.UUID  `json:"ownerId,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	Keys      []DiaryKey `json:"keys,omitempty"`
}

type DiaryKey struct {
	DiaryId      uuid.UUID `json:"diaryId,omitempty"`
	UserId       uuid.UUID `json:"userId,omitempty"`
	EncryptedKey string    `json:"encryptedKey,omitempty"`
}

type DiaryUseCase interface {
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}

type DiaryRepository interface {
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]Diary, error)
}
