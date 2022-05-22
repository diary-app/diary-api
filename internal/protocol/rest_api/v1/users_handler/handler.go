package users_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetMe() gin.HandlerFunc
	GetUserById() gin.HandlerFunc
	GetUserByName() gin.HandlerFunc
}

func New(uc usecase.UsersUseCase) Handler {
	return &handler{
		uc: uc,
	}
}

type handler struct {
	uc usecase.UsersUseCase
}
