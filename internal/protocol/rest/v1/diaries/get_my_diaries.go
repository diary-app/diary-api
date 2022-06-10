package diaries

import (
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetMyDiaries() gin.HandlerFunc {
	return func(c *gin.Context) {
		diaries, err := h.uc.GetDiariesByUser(c)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, usecase.DiaryListResponse{Items: mapList(diaries)})
	}
}

func mapList(diaries []usecase.Diary) []usecase.DiaryResponse {
	list := make([]usecase.DiaryResponse, len(diaries))
	for i, d := range diaries {
		list[i] = mapToDiaryResponse(&d)
	}
	return list
}

func mapToDiaryResponse(d *usecase.Diary) usecase.DiaryResponse {
	return usecase.DiaryResponse{
		ID:           d.ID,
		Name:         d.Name,
		OwnerID:      d.OwnerID,
		EncryptedKey: d.Keys[0].EncryptedKey,
	}
}
