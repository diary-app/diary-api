package users

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := utils.ParseUUIDFromPath(c, "id")
		if !ok {
			return
		}

		user, err := h.uc.GetUserByID(c, id)
		if err != nil {
			if err == usecase.ErrUserNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound,
					common.ErrorResponse{Message: fmt.Sprintf("user with id '%s' was not found", id)})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
