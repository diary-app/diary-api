package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := utils.ParseUUIDFromPath(c, "id")
		if !ok {
			return
		}

		diaryEntry, err := h.uc.GetByID(c, id)
		if err != nil {
			noAccessErr, ok := err.(*usecase.NoReadAccessToDiaryEntryError)
			if ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: noAccessErr.Error()})
				return
			}
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, mapToEntryResponse(diaryEntry))
	}
}
