//go:generate mockgen -source ${GOFILE} -package ${GOPACKAGE}_test -destination mocks_test.go
package usecase

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type storage interface {
	CreateUser(
		ctx context.Context, user *usecase.FullUser, diary *usecase.Diary) (*usecase.FullUser, *usecase.Diary, error)
	GetUserByName(ctx context.Context, username string) (*usecase.FullUser, error)
}

type tokenService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	RefreshToken(token string) (string, error)
}
