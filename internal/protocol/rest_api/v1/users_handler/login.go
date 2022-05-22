package users_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (h *handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.LoginRequest{}
		if err := c.BindJSON(request); err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "request body was invalid"})
			return
		}

		authResult, err := h.uc.Login(c.Request.Context(), request)
		if err != nil {
			if err == usecase.UserNotFoundError || err == bcrypt.ErrMismatchedHashAndPassword {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "incorrect user or password"})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, authResult)
	}
}
