package auth

import (
	"errors"
	"fmt"
)

type ErrEarlyForTokenRefresh struct {
	SecondsBeforeExpire float64
}

func NewEarlyForTokenRefreshError(secondsUntilExpire float64) ErrEarlyForTokenRefresh {
	return ErrEarlyForTokenRefresh{
		secondsUntilExpire,
	}
}

func (e ErrEarlyForTokenRefresh) Error() string {
	return fmt.Sprintf(
		"too early: token refresh is possible within %v minutes, but token expires in %v seconds",
		MinutesBeforeExpireToRefresh, e.SecondsBeforeExpire)
}

var (
	ErrTokenInvalid           = errors.New("jwt token was invalid")
	ErrUserIdNotFoundInCtx    = errors.New("user id not found in request context")
	ErrAuthTokenNotFoundInCtx = errors.New("auth token not found in request context")
)
