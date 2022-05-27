package auth

import (
	"diary-api/internal/config"
	"github.com/benbjohnson/clock"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

const (
	MinutesBeforeExpireToRefresh = 5
	TokenLifespanMinutes         = 60
)

type Claims struct {
	UserId uuid.UUID
	jwt.StandardClaims
}

type TokenService interface {
	GenerateToken(userId uuid.UUID) (string, error)
	ValidateToken(token string) (*Claims, error)
	RefreshToken(token string) (string, error)
}

func NewAuthService(cfg *config.AuthConfig, c clock.Clock) TokenService {
	return &tokensManager{
		jwtKey: []byte(cfg.JwtKey),
		clock:  c,
	}
}

type tokensManager struct {
	jwtKey []byte
	clock  clock.Clock
}

func (t *tokensManager) GenerateToken(userId uuid.UUID) (string, error) {
	exp := t.clock.Now().Add(TokenLifespanMinutes * time.Minute)
	claims := &Claims{
		UserId: userId,
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

func (t *tokensManager) ValidateToken(tokenString string) (*Claims, error) {
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

func (t *tokensManager) RefreshToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, t.getJwtKey)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrTokenInvalid
	}

	untilExpire := time.Unix(claims.ExpiresAt, 0).Sub(t.clock.Now())
	if untilExpire > MinutesBeforeExpireToRefresh*time.Minute {
		return "", NewEarlyForTokenRefreshError(untilExpire.Seconds())
	}

	return t.GenerateToken(claims.UserId)
}

func (t *tokensManager) getJwtKey(_ *jwt.Token) (interface{}, error) {
	return t.jwtKey, nil
}
