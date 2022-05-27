package users_handler

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *handler) GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'id' in path should not be empty"})
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'id' in path was not a valid UUID"})
			return
		}
		user, err := h.uc.GetUserById(c, id)
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
