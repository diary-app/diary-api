package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) PatchEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := utils.ParseUUIDFromPath(c, "id")
		if !ok {
			return
		}

		request := &usecase.UpdateDiaryEntryRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			utils.RespondInvalidBodyJSONWithError(c, err)
			return
		}

		err := h.uc.Update(c, id, request)
		if err != nil {
			if alienErr, ok := err.(*usecase.AlienEntryBlocksError); ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: alienErr.Error()})
			} else if diaryAccessErr, ok := err.(*usecase.NoAccessToDiaryError); ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: diaryAccessErr.Error()})
			} else if entryAccessErr, ok := err.(*usecase.NoAccessToDiaryEntryError); ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: entryAccessErr.Error()})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
