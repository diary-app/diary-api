package diary_entries

import (
	"diary-api/internal/protocol/rest/utils"
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
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
