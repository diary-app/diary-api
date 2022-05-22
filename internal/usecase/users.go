package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type FullUser struct {
	Id                            uuid.UUID `db:"id"`
	Username                      string    `db:"username"`
	PasswordHash                  []byte    `db:"password_hash"`
	SaltForKeys                   []byte    `db:"salt_for_keys"`
	PublicKeyForSharing           string    `db:"public_key_for_sharing"`
	EncryptedPrivateKeyForSharing string    `db:"encrypted_private_key_for_sharing"`
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
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type UsersUseCase interface {
	Register(ctx context.Context, request *RegisterRequest) (*RegistrationResult, error)
	Login(ctx context.Context, request *LoginRequest) (*AuthResult, error)
	GetFullUser(ctx context.Context, userId uuid.UUID) (*FullUser, error)
	GetUserByName(ctx context.Context, username string) (*ShortUser, error)
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *FullUser) (*FullUser, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*FullUser, error)
	GetUserByName(ctx context.Context, username string) (*FullUser, error)
}

// Errors

var (
	UserNotFoundError = errors.New("user with given name was not found")
)

type UsernameTakenError struct {
	Username string
}

func (u UsernameTakenError) Error() string {
	return fmt.Sprintf("username '%v' is already taken", u.Username)
}
