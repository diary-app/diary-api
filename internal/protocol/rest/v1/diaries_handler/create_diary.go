package diaries_handler

import (
	"diary-api/internal/auth"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) CreateDiary() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.CreateDiaryRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			utils.RespondInvalidBodyJSON(c)
			return
		}

		userId := auth.MustGetUserId(c)
		diary, err := h.uc.CreateDiary(c, userId, request)
		if err != nil {
			if err == usecase.ErrDuplicateDiaryName {
				c.JSON(http.StatusConflict,
					gin.H{"message": fmt.Sprintf("you already have diary with name '%v'", request.Name)})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, diary)
	}
}
