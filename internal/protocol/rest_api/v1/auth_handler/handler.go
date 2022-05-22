package auth_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
}

type handler struct {
	uc usecase.UsersUseCase
}

func New(uc usecase.UsersUseCase) Handler {
	return &handler{uc: uc}
}
