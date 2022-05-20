package diary_entries_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetList() gin.HandlerFunc
	Download() gin.HandlerFunc
	Create() gin.HandlerFunc
	Upload() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
}

func New(uc usecase.DiaryEntriesUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.DiaryEntriesUseCase
}
