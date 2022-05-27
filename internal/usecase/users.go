package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type FullUser struct {
	Id                            uuid.UUID `db:"id" json:"id"`
	Username                      string    `db:"username" json:"username"`
	PasswordHash                  []byte    `db:"password_hash" json:"-"`
	SaltForKeys                   []byte    `db:"salt_for_keys" json:"saltForKeys"`
	PublicKeyForSharing           string    `db:"public_key_for_sharing" json:"publicKeyForSharing"`
	EncryptedPrivateKeyForSharing string    `db:"encrypted_private_key_for_sharing" json:"encryptedPrivateKeyForSharing"`
}

func (u *FullUser) String() string {
	return fmt.Sprintf("%s (%s)", u.Username, u.Id.String())
}

type ShortUser struct {
	Id                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	PublicKeyForSharing string    `json:"public_key_for_sharing"`
}

type RegisterRequest struct {
	Username                      string `json:"username" binding:"required"`
	Password                      string `json:"password" binding:"required"`
	PublicKeyForSharing           string `json:"publicKeyForSharing" binding:"required"`
	EncryptedPrivateKeyForSharing string `json:"encryptedPrivateKeyForSharing" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrationResult struct {
	Token       string `json:"token"`
	SaltForKeys string `json:"saltForKeys"`
}

type AuthResult struct {
	Token string `json:"token"`
}

type UsersUseCase interface {
	Register(ctx context.Context, request *RegisterRequest) (*RegistrationResult, error)
	Login(ctx context.Context, request *LoginRequest) (*AuthResult, error)
	RefreshToken(ctx context.Context, token string) (*AuthResult, error)
	GetFullUser(ctx context.Context, userId uuid.UUID) (*FullUser, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*ShortUser, error)
	GetUserByName(ctx context.Context, username string) (*ShortUser, error)
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *FullUser) (*FullUser, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*FullUser, error)
	GetUserByName(ctx context.Context, username string) (*FullUser, error)
}

// Errors

var (
	ErrUserNotFound = errors.New("user with given name was not found")
)

type ErrUsernameTaken struct {
	Username string
}

func (u ErrUsernameTaken) Error() string {
	return fmt.Sprintf("username '%v' is already taken", u.Username)
}
