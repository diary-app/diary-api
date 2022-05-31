package sharing_tasks

import (
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetSharingTasks() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.MustGetUserID(c)
		tasks, err := h.uc.GetSharingTasks(c, userID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, usecase.SharingTasksListResponse{Items: tasks})
	}
}
