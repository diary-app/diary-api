package sharing_tasks

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create() gin.HandlerFunc
	GetSharingTasks() gin.HandlerFunc
	AcceptSharedDiary() gin.HandlerFunc
}

func New(uc usecase.SharingTasksUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.SharingTasksUseCase
}
