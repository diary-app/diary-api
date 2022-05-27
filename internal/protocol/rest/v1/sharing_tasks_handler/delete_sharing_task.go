package sharing_tasks_handler

import (
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const (
	PathDiaryIdKey = "diaryId"
)

func (h *handler) DeleteById() gin.HandlerFunc {
	return func(c *gin.Context) {
		diaryIdStr := c.Param(PathDiaryIdKey)
		if diaryIdStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"message": fmt.Sprintf("path should contain diary id (/%v=:diaryId)", PathDiaryIdKey)})
			return
		}

		diaryId, err := uuid.Parse(diaryIdStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"message": fmt.Sprintf("%v should be a valid UUID", PathDiaryIdKey)})
			return
		}

		userId := auth.MustGetUserId(c)
		err = h.uc.DeleteSharingTask(c, diaryId, userId)
		if err != nil && err != usecase.ErrCommonNotFound {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
