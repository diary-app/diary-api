package diary_entries

import (
	"diary-api/internal/protocol/rest/common"
	"diary-api/internal/protocol/rest/utils"
	"diary-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.CreateDiaryEntryRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			utils.RespondInvalidBodyJSONWithError(c, err)
			return
		}

		entry, err := h.uc.Create(c, *request)
		if err != nil {
			accessErr, ok := err.(*usecase.NoAccessToDiaryError)
			if ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ErrorResponse{Message: accessErr.Error()})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, mapToEntryResponse(entry))
	}
}
