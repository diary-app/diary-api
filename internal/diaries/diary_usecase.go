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

func (uc *diaryUseCase) CreateDiary(ctx context.Context, userId uuid.UUID,
	req *usecase.CreateDiaryRequest) (*usecase.Diary, error) {
	diaryKeys := make([]usecase.DiaryKey, 1)
	diaryKeys[0] = usecase.DiaryKey{
		UserId:       userId,
		EncryptedKey: req.EncryptedKey,
	}
	diary := &usecase.Diary{
		Name:    req.Name,
		OwnerId: userId,
		Keys:    diaryKeys,
	}

	diary, err := uc.repo.CreateDiary(ctx, diary)
	if err != nil {
		return nil, err
	}

	return diary, nil
}

func (uc *diaryUseCase) GetDiariesByUser(ctx context.Context, userId uuid.UUID) ([]usecase.Diary, error) {
	return uc.repo.GetDiariesByUser(ctx, userId)
}
