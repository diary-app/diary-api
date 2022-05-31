package usecase

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
)

type UseCase struct {
	s storage
}

func New(s storage) *UseCase {
	return &UseCase{
		s: s,
	}
}

func (uc *UseCase) GetDiariesByUser(ctx context.Context) ([]usecase.Diary, error) {
	userID := auth.MustGetUserID(ctx)
	return uc.s.GetDiariesByUser(ctx, userID)
}
