package auth

import (
	"errors"
	"fmt"
)

type EarlyForTokenRefreshError struct {
	SecondsBeforeExpire float64
}

func (e EarlyForTokenRefreshError) Error() string {
	return fmt.Sprintf(
		"too early: token refresh is possible within %v seconds of expiration, but token expires in %v seconds",
		SecondsBeforeExpireToRefresh.Seconds(), e.SecondsBeforeExpire)
}

var (
	ErrTokenInvalid           = errors.New("jwt token was invalid")
	ErrUserIDNotFoundInCtx    = errors.New("user id not found in request context")
	ErrAuthTokenNotFoundInCtx = errors.New("auth token not found in request context")
)
