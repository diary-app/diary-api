package diaries_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	GetMyDiaries() gin.HandlerFunc
}

func New(usecase usecase.DiaryUseCase) Handler {
	return &handler{
		usecase: usecase,
	}
}

type handler struct {
	usecase usecase.DiaryUseCase
	l       log.Logger
}
