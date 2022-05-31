package diaries

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
		uc: usecase,
	}
}

type handler struct {
	uc usecase.DiaryUseCase
	l  log.Logger
}
