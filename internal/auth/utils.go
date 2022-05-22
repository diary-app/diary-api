package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MustGetUserId(c *gin.Context) uuid.UUID {
	val, ok := c.Get(UserIdKey)
	if !ok {
		panic(ErrUserIdNotFoundInCtx)
	}

	id := val.(uuid.UUID)
	return id
}
