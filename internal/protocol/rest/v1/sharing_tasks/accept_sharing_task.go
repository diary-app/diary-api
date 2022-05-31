package sharing_tasks

import (
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const (
	PathDiaryIDKey = "diaryID"
)

func (h *handler) AcceptByDiaryID() gin.HandlerFunc {
	return func(c *gin.Context) {
		diaryIDStr := c.Param(PathDiaryIDKey)
		if diaryIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"message": fmt.Sprintf("path should contain diary id (/%v=:diaryID)", PathDiaryIDKey)})
			return
		}

		diaryID, err := uuid.Parse(diaryIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"message": fmt.Sprintf("%v should be a valid UUID", PathDiaryIDKey)})
			return
		}

		userID := auth.MustGetUserID(c)
		err = h.uc.DeleteSharingTask(c, diaryID, userID)
		if err != nil && err != usecase.ErrCommonNotFound {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
