package diary_entries

import (
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &usecase.CreateDiaryEntryRequest{}
		if err := c.ShouldBindJSON(request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid request body: %v", err)})
			return
		}

		entry, err := h.uc.Create(c, *request)
		if err != nil {
			accessErr, ok := err.(*usecase.NoAccessToDiaryError)
			if ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": accessErr.Error()})
			} else {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}

		c.JSON(http.StatusOK, mapToEntryResponse(entry))
	}
}
