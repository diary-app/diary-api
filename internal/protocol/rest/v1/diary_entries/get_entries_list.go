package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (h *handler) GetEntries() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := usecase.GetDiaryEntriesParamsDto{}
		if err := c.ShouldBind(&req); err != nil {
			utils.RespondInvalidQueryParams(c)
			return
		}
		params := usecase.GetDiaryEntriesParams{}
		if req.DiaryIDStr != nil {
			if parsed, err := uuid.Parse(*req.DiaryIDStr); err != nil {
				utils.RespondInvalidQueryParams(c)
				return
			} else {
				params.DiaryID = &parsed
			}
		}
		if req.DateStr != nil {
			if parsed, err := time.Parse("2006-01-02", *req.DateStr); err != nil {
				utils.RespondInvalidQueryParams(c)
				return
			} else {
				dateOnly := common.DateOnly(parsed)
				params.Date = &dateOnly
			}
		}
		entries, err := h.uc.GetEntries(c, params)
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
