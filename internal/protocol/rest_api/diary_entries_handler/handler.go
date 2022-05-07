package diary_entries_handler

import "github.com/gin-gonic/gin"

type Handler interface {
	GetList() gin.HandlerFunc
	Download() gin.HandlerFunc
	Create() gin.HandlerFunc
	Upload() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
}

func New() Handler {
	return &handler{}
}

type handler struct {
}
