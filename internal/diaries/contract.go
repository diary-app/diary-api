//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go

package diaries

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type storage interface {
	CreateDiary(ctx context.Context, diary *usecase.Diary) (*usecase.Diary, error)
	GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]usecase.Diary, error)
}
