package users

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetUserByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("name")
		if username == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'name' in path should not be empty"})
			return
		}

		user, err := h.uc.GetUserByName(c, username)
		if err != nil {
			if err == usecase.ErrUserNotFound {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
