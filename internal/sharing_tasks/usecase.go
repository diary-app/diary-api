package sharing_tasks

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type sharingTasksUseCase struct {
	repo usecase.SharingTasksRepository
}

func (uc *sharingTasksUseCase) CreateSharingTask(ctx context.Context, req *usecase.NewSharingTaskRequest) error {
	st := &usecase.SharingTask{
		DiaryID:           req.DiaryID,
		ReceiverUserID:    req.ReceiverUserID,
		EncryptedDiaryKey: req.EncryptedDiaryKey,
		SharedAt:          req.SharedAt,
	}
	return uc.repo.CreateSharingTask(ctx, st)
}

func (uc *sharingTasksUseCase) GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]usecase.SharingTask, error) {
	return uc.repo.GetSharingTasks(ctx, userID)
}

func (uc *sharingTasksUseCase) DeleteSharingTask(ctx context.Context, diaryID uuid.UUID, receiverID uuid.UUID) error {
	return uc.repo.DeleteSharingTask(ctx, diaryID, receiverID)
}

func NewUseCase(repo usecase.SharingTasksRepository) usecase.SharingTasksUseCase {
	return &sharingTasksUseCase{repo: repo}
}
