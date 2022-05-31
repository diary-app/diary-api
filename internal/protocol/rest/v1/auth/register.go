package auth

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.RegisterRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			utils.RespondInvalidBodyJSON(c)
			return
		}

		authResult, err := h.useCase.Register(c, request)
		if err != nil {
			usernameTakenError, ok := err.(usecase.ErrUsernameTaken)
			if ok {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": usernameTakenError.Error()})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, authResult)
	}
}
