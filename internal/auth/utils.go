package auth

import (
	"context"
	"github.com/google/uuid"
)

func getUserId(ctx context.Context) (uuid.UUID, error) {
	idStr := ctx.Value(UserIdKey)
	if idStr == "" {
		return uuid.UUID{}, ErrUserIdNotFoundInCtx
	}

	id := idStr.(uuid.UUID)
	return id, nil
}

func MustGetUserId(ctx context.Context) uuid.UUID {
	id, err := getUserId(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

func MustGetAuthToken(ctx context.Context) string {
	authToken, ok := ctx.Value(AuthToken).(string)
	if !ok || authToken == "" {
		panic(ErrAuthTokenNotFoundInCtx)
	}

	return authToken
}
