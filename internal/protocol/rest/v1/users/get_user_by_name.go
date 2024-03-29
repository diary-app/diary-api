package users

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetUserByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("name")
		if username == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: "'name' in path should not be empty"})
			return
		}

		user, err := h.uc.GetUserByName(c, username)
		if err != nil {
			if err == usecase.ErrUserNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound,
					common.ErrorResponse{Message: fmt.Sprintf("user with name '%s' not found", username)})
				return
			}

			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
