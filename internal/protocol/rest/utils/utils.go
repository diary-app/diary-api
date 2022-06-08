package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func RespondInvalidBodyJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
}

func ParseUUIDFromPath(c *gin.Context, key string) (id uuid.UUID, parsed bool) {
	idStr := c.Param(key)
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'id' in path should not be empty"})
		return uuid.UUID{}, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "'id' in path should be a valid UUID"})
		return uuid.UUID{}, false
	}

	return id, true
}
