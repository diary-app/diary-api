package diary_entries_handler

import (
	"diary-api/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &usecase.CreateDiaryEntryRequest{}
		if err := ctx.BindJSON(request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid request body: %v", err)})
			return
		}

		entry, err := h.uc.Create(ctx.Request.Context(), *request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, entry)
	}
}
