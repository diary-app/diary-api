package diary_entries

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetEntries() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Create() gin.HandlerFunc
	PatchEntry() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

func New(uc usecase.DiaryEntriesUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.DiaryEntriesUseCase
}
