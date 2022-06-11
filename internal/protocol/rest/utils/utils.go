package utils

import (
	"diary-api/internal/protocol/rest/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func RespondInvalidQueryParams(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid query params"})
}

func RespondInvalidBodyJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid request body"})
}

func RespondInvalidBodyJSONWithError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest,
		common.ErrorResponse{Message: fmt.Sprintf("invalid requuest body: %s", err.Error())})
}

func ParseUUIDFromPath(c *gin.Context, key string) (id uuid.UUID, parsed bool) {
	idStr := c.Param(key)
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: "'id' in path should not be empty"})
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Message: "'id' in path should be a valid UUID"})
		return uuid.UUID{}, false
	}

	return id, true
}
