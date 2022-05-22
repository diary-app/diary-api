package rest_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondInvalidBodyJSON(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
}
