package auth

import (
	"context"
	"github.com/google/uuid"
)

func getUserID(ctx context.Context) (uuid.UUID, error) {
	idStr := ctx.Value(UserIDKey)
	if idStr == "" {
		return uuid.UUID{}, ErrUserIDNotFoundInCtx
	}

	id := idStr.(uuid.UUID)
	return id, nil
}

func MustGetUserID(ctx context.Context) uuid.UUID {
	id, err := getUserID(ctx)
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
