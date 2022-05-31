package users

import (
	"diary-api/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.MustGetUserID(c)
		user, err := h.uc.GetFullUser(c, userID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
