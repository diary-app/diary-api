package diary_entries

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "diary id 'id' in path is missing"})
			return
		}

	}
}
