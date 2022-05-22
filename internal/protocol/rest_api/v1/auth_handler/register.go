package auth_handler

import (
	"diary-api/internal/protocol/rest_api/rest_utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.RegisterRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			rest_utils.RespondInvalidBodyJSON(c)
			return
		}

		authResult, err := h.uc.Register(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			usernameTakenError, ok := err.(usecase.ErrUsernameTaken)
			if ok {
				c.JSON(http.StatusConflict, gin.H{"message": usernameTakenError.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, authResult)
	}
}
