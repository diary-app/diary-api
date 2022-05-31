package auth

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (h *handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.LoginRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			utils.RespondInvalidBodyJSON(c)
			return
		}

		authResult, err := h.useCase.Login(c, request)
		if err != nil {
			if err == usecase.ErrUserNotFound || err == bcrypt.ErrMismatchedHashAndPassword {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "incorrect user or password"})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, authResult)
	}
}
