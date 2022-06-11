package sharing_tasks

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type sharingTasksUseCase struct {
	repo usecase.SharingTasksRepository
}

func (uc *sharingTasksUseCase) CreateSharingTask(
	ctx context.Context,
	req *usecase.CreateSharingTaskRequest,
) (*usecase.SharingTask, error) {
	return uc.repo.CreateSharingTask(ctx, req)
}

func (uc *sharingTasksUseCase) GetSharingTasks(ctx context.Context, userID uuid.UUID) ([]usecase.SharingTask, error) {
	return uc.repo.GetSharingTasks(ctx, userID)
}

func (uc *sharingTasksUseCase) AcceptSharingTask(ctx context.Context, req *usecase.AcceptSharingTaskRequest) error {
	return uc.repo.AcceptSharingTask(ctx, req)
}

func NewUseCase(repo usecase.SharingTasksRepository) usecase.SharingTasksUseCase {
	return &sharingTasksUseCase{repo: repo}
}
