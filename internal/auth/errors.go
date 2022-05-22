package auth

import (
	"errors"
	"fmt"
)

type EarlyForTokenRefreshError struct {
	SecondsBeforeExpire float64
}

func NewEarlyForTokenRefreshError(secondsUntilExpire float64) EarlyForTokenRefreshError {
	return EarlyForTokenRefreshError{
		secondsUntilExpire,
	}
}

func (e EarlyForTokenRefreshError) Error() string {
	return fmt.Sprintf("too early for token refresh, token expires in %v seconds", e.SecondsBeforeExpire)
}

var (
	TokenInvalidError   = errors.New("jwt token was invalid")
	UserIdNotFoundError = errors.New("user id not found in request context")
)
