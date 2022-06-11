package sharing_tasks

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &usecase.NewSharingTaskRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			utils.RespondInvalidBodyJSONWithError(c, err)
			return
		}

		task, err := h.uc.CreateSharingTask(c, req)
		if err == nil {
			c.JSON(http.StatusOK, usecase.NewSharingTaskResponse{DiaryID: task.DiaryID})
			return
		}

		if err == usecase.ErrReceiverUserNotFound {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		} else if accessErr, ok := err.(*usecase.NoAccessToDiaryEntryError); ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: accessErr.Error()})
		} else if err == usecase.ErrUserAlreadyHasAccessToDiary {
			c.AbortWithStatusJSON(http.StatusConflict, common.ErrorResponse{Message: "user already has access to the diary"})
		} else if err == usecase.ErrUserAlreadyHasTaskForSameDiary {
			c.AbortWithStatusJSON(http.StatusConflict, common.ErrorResponse{Message: "user already has sharing task for the diary"})
		} else if blocksErr, ok := err.(*usecase.BadUpdatedBlocksError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: blocksErr.Error()})
		} else if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
