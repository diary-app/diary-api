package diaries_handler

import (
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) CreateDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.CreateDiaryRequest{}
		if err := c.BindJSON(request); err != nil {
			_ = c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
			return
		}

		userId := auth.MustGetUserId(c)
		diary, err := h.uc.CreateDiary(c.Request.Context(), userId, request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, diary)
	}
}
