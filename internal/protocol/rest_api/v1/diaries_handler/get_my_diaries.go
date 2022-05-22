package diaries_handler

import (
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type DiariesResponse struct {
	Items []usecase.Diary `json:"items"`
}

func (h *handler) GetMyDiaries() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId uuid.UUID
		// TODO userId from jwt
		userId = uuid.MustParse("de0a2016-6df3-437f-88cb-bec859f3c53f")
		diaries, err := h.uc.GetDiariesByUser(c.Request.Context(), userId)
		if err != nil {
			// TODO log
			_ = c.Error(fmt.Errorf("GetMyDiaries - DiaryUseCase - GetDiariesByUser: %v", err))
			return
		}

		response := DiariesResponse{
			Items: diaries,
		}

		c.JSON(http.StatusOK, response)
	}
}
