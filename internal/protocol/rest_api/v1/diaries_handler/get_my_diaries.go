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

func (c *handler) GetMyDiaries() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userId uuid.UUID
		// TODO userId from jwt
		userId = uuid.MustParse("de0a2016-6df3-437f-88cb-bec859f3c53f")
		diaries, err := c.usecase.GetDiariesByUser(ctx.Request.Context(), userId)
		if err != nil {
			// TODO log
			_ = ctx.Error(fmt.Errorf("GetMyDiaries - DiaryUseCase - GetDiariesByUser: %v", err))
			return
		}

		response := DiariesResponse{
			Items: diaries,
		}

		ctx.JSON(http.StatusOK, response)
	}
}
