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
		DiaryId:           req.DiaryId,
		ReceiverUserId:    req.ReceiverUserId,
		EncryptedDiaryKey: req.EncryptedDiaryKey,
		SharedAt:          req.SharedAt,
	}
	return uc.repo.CreateSharingTask(ctx, st)
}

func (uc *sharingTasksUseCase) GetSharingTasks(ctx context.Context, userId uuid.UUID) ([]usecase.SharingTask, error) {
	return uc.repo.GetSharingTasks(ctx, userId)
}

func (uc *sharingTasksUseCase) DeleteSharingTask(ctx context.Context, diaryId uuid.UUID, receiverId uuid.UUID) error {
	return uc.repo.DeleteSharingTask(ctx, diaryId, receiverId)
}

func NewUseCase(repo usecase.SharingTasksRepository) usecase.SharingTasksUseCase {
	return &sharingTasksUseCase{repo: repo}
}
