package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type FullUser struct {
	ID                            uuid.UUID `db:"id" json:"id"`
	Username                      string    `db:"username" json:"username"`
	PasswordHash                  []byte    `db:"password_hash" json:"-"`
	MasterKeySalt                 []byte    `db:"salt_for_keys" json:"masterKeySalt"`
	PublicKeyForSharing           string    `db:"public_key_for_sharing" json:"publicKeyForSharing"`
	EncryptedPrivateKeyForSharing []byte    `db:"encrypted_private_key_for_sharing" json:"encryptedPrivateKeyForSharing"`
}

func (u *FullUser) String() string {
	return fmt.Sprintf("%s (%s)", u.Username, u.ID.String())
}

type ShortUser struct {
	ID                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	PublicKeyForSharing string    `json:"publicKeyForSharing"`
}

type RegisterRequest struct {
	Username                      string `json:"username" binding:"required"`
	Password                      string `json:"password" binding:"required"`
	MasterKeySalt                 []byte `json:"masterKeySalt" binding:"required"`
	PublicKeyForSharing           string `json:"publicKeyForSharing" binding:"required"`
	EncryptedPrivateKeyForSharing []byte `json:"encryptedPrivateKeyForSharing" binding:"required"`
	EncryptedDiaryKey             []byte `json:"encryptedDiaryKey" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrationResult struct {
	Token   string    `json:"token"`
	DiaryID uuid.UUID `json:"diaryId"`
}

type AuthResult struct {
	Token string `json:"token"`
}

type UsersUseCase interface {
	GetFullUser(ctx context.Context, userID uuid.UUID) (*FullUser, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*ShortUser, error)
	GetUserByName(ctx context.Context, username string) (*ShortUser, error)
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user *FullUser, diary *Diary) (*FullUser, *Diary, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*FullUser, error)
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
