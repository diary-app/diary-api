package diary_entries

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetEntriesList() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := usecase.GetDiaryEntriesParams{}
		if err := c.ShouldBindUri(&req); err != nil {
			utils.RespondInvalidBodyJSON(c)
			return
		}
		entries, err := h.uc.GetEntries(c, req)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, usecase.GetDiaryEntriesResponse{Items: entries})
	}
}
