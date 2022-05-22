package diaries_handler

import (
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DiariesResponse struct {
	Items []usecase.Diary `json:"items"`
}

func (h *handler) GetMyDiaries() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := auth.MustGetUserId(c)
		diaries, err := h.uc.GetDiariesByUser(c.Request.Context(), userId)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, DiariesResponse{Items: diaries})
	}
}
