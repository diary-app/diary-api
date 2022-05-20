package sharing_tasks_handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Create() gin.HandlerFunc
	GetAllMine() gin.HandlerFunc
	DeleteById() gin.HandlerFunc
}

func New() Handler {
	return &handler{}
}

type handler struct {
}
