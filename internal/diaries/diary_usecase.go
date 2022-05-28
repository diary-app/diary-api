package diaries

import (
	"context"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type UseCase struct {
	s storage
}

func NewDiaryUseCase(s storage) *UseCase {
	return &UseCase{
		s: s,
	}
}

func (uc *UseCase) CreateDiary(
	ctx context.Context,
	userId uuid.UUID,
	req *usecase.CreateDiaryRequest,
) (*usecase.Diary, error) {
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

	diary, err := uc.s.CreateDiary(ctx, diary)
	if err != nil {
		return nil, err
	}

	return diary, nil
}

func (uc *UseCase) GetDiariesByUser(ctx context.Context) ([]usecase.Diary, error) {
	userId := auth.MustGetUserId(ctx)
	return uc.s.GetDiariesByUser(ctx, userId)
}
