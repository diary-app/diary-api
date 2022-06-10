package sharing_tasks

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) AcceptSharedDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &usecase.AcceptSharingTaskRequest{}
		if err := c.ShouldBindJSON(req); err != nil {
			utils.RespondInvalidBodyJSON(c)
			return
		}

		err := h.uc.AcceptSharingTask(c, req)
		if err != nil && err != usecase.ErrCommonNotFound {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
