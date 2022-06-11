package auth

import (
	"diary-api/internal/config"
	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

const (
	SecondsBeforeExpireToRefresh = 5 * time.Minute
)

type Claims struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type TokenService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (*Claims, error)
	RefreshToken(token string) (string, error)
}

func NewTokenService(cfg *config.AuthConfig, c clock.Clock) TokenService {
	return &tokenService{
		jwtKey: []byte(cfg.JwtKey),
		jwtTtl: time.Duration(cfg.JwtTtlMinutes) * time.Minute,
		clock:  c,
	}
}

type tokenService struct {
	jwtKey []byte
	jwtTtl time.Duration
	clock  clock.Clock
}

func (t *tokenService) GenerateToken(userID uuid.UUID) (string, error) {
	exp := t.clock.Now().Add(t.jwtTtl)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *tokenService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, t.getJwtKey)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (t *tokenService) RefreshToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, t.getJwtKey)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrTokenInvalid
	}

	untilExpire := time.Unix(claims.ExpiresAt, 0).Sub(t.clock.Now())
	if untilExpire > SecondsBeforeExpireToRefresh {
		return "", EarlyForTokenRefreshError{untilExpire.Seconds()}
	}

	return t.GenerateToken(claims.UserID)
}

func (t *tokenService) getJwtKey(_ *jwt.Token) (interface{}, error) {
	return t.jwtKey, nil
}
