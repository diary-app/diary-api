package diaries

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type diaryUseCase struct {
	repo usecase.DiaryRepository
}

func NewDiaryUseCase(repo usecase.DiaryRepository) usecase.DiaryUseCase {
	return &diaryUseCase{
		repo: repo,
	}
}

func (d *diaryUseCase) GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]usecase.Diary, error) {
	return d.repo.GetDiariesByUser(ctx, userId)
}
