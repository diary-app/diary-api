package users_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	GetMe() gin.HandlerFunc
	GetUser() gin.HandlerFunc
}

func New(uc usecase.UsersUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.UsersUseCase
}
