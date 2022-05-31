package auth

import (
	"diary-api/internal/auth"
	"diary-api/internal/protocol/rest/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.MustGetAuthToken(c)
		authResult, err := h.useCase.RefreshToken(c, token)
		if err == nil {
			c.JSON(http.StatusOK, authResult)
			return
		}
		if earlyErr, ok := err.(*auth.ErrEarlyForTokenRefresh); ok {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{Message: earlyErr.Error()})
			return
		}

		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
}
