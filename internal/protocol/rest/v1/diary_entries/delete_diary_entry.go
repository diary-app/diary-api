package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := utils.ParseUUIDFromPath(c, "id")
		if !ok {
			return
		}

		err := h.uc.Delete(c, id)
		if err != nil {
			if accessErr, ok := err.(*usecase.NoWriteAccessToDiaryEntryError); ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: accessErr.Error()})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.Status(http.StatusNoContent)
	}
}
