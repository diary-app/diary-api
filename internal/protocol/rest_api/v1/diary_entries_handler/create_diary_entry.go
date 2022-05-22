package diary_entries_handler

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

		entry, err := h.uc.Create(c.Request.Context(), *request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}
