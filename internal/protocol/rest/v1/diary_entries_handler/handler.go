package diary_entries_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetEntriesList() gin.HandlerFunc
	GetById() gin.HandlerFunc
	Create() gin.HandlerFunc
	PatchEntry() gin.HandlerFunc
	UpdateContents() gin.HandlerFunc
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
