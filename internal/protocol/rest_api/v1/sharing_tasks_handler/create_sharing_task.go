package sharing_tasks_handler

import (
	"diary-api/internal/protocol/rest_api/rest_utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &usecase.NewSharingTaskRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			rest_utils.RespondInvalidBodyJSON(c)
			return
		}

		err := h.uc.CreateSharingTask(c.Request.Context(), req)
		if err == usecase.ErrUserAlreadyHasAccessToDiary {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "user already has access to the diary"})
			return
		}
		if err == usecase.ErrUserAlreadyHasTaskForSameDiary {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "user already has sharing task for the diary"})
			return
		}
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
