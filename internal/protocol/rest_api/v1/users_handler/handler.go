package users_handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	GetMe() gin.HandlerFunc
	GetUser() gin.HandlerFunc
}

func New() Handler {
	return &handler{}
}

type handler struct {
}
