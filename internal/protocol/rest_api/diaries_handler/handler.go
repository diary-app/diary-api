package diaries_handler

import "github.com/gin-gonic/gin"

type Handler interface {
	GetMyDiaries() gin.HandlerFunc
	GetDiaryEntries() gin.HandlerFunc
}

func New() Handler {
	return &handler{}
}

type handler struct {
}
