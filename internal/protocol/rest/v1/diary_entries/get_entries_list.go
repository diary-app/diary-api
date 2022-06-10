package diary_entries

import (
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetEntries() gin.HandlerFunc {
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

		shortDiaries := mapShortEntryResponseList(entries)
		c.JSON(http.StatusOK, DiaryEntriesResponse{Items: shortDiaries})
	}
}

type DiaryEntriesResponse struct {
	Items []usecase.ShortDiaryEntryResponse `json:"items"`
}

func mapShortEntryResponseList(entries []usecase.DiaryEntry) []usecase.ShortDiaryEntryResponse {
	response := make([]usecase.ShortDiaryEntryResponse, len(entries))
	for i, e := range entries {
		response[i] = mapToShortEntry(&e)
	}
	return response
}
