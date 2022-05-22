package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MustGetUserId(ctx *gin.Context) uuid.UUID {
	val, ok := ctx.Get(UserIdKey)
	if !ok {
		panic(UserIdNotFoundError)
	}

	id := val.(uuid.UUID)
	return id
}
