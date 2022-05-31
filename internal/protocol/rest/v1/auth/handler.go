package auth

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register() gin.HandlerFunc
	Login() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
}

type handler struct {
	useCase usecase.AuthUseCase
}

func New(uc usecase.AuthUseCase) Handler {
	return &handler{useCase: uc}
}
