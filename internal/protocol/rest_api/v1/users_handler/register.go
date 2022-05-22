package users_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.RegisterRequest{}
		if err := c.BindJSON(request); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "request body was invalid"})
			return
		}

		authResult, err := h.uc.Register(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			usernameTakenError, ok := err.(usecase.UsernameTakenError)
			if ok {
				c.JSON(http.StatusConflict, gin.H{"message": usernameTakenError.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, authResult)
	}
}
