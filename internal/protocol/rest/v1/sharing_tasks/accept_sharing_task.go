package sharing_tasks

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) AcceptSharedDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &usecase.AcceptSharingTaskRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			utils.RespondInvalidBodyJSONWithError(c, err)
			return
		}

		err := h.uc.AcceptSharingTask(c, req)
		if err != nil {
			if err == usecase.ErrCommonNotFound {
				c.AbortWithStatusJSON(http.StatusNotFound,
					common.ErrorResponse{Message: fmt.Sprintf("current user does not have task for diary %v", req.DiaryID)})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
