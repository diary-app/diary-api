package diary_entries

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"fmt"
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid request body: %v", err)})
			return
		}

		err := h.uc.Update(c, id, request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
