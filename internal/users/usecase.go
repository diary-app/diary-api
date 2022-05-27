package users

import (
	"context"
	"database/sql"
	"diary-api/internal/auth"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

const (
	SaltForKeysSize = 16
)

type usersUseCase struct {
	tokensService auth.TokenService
	usersRepo     usecase.UsersRepository
	diariesRepo   usecase.DiaryRepository
}

func NewUseCase(t auth.TokenService, r usecase.UsersRepository) usecase.UsersUseCase {
	return &usersUseCase{
		tokensService: t,
		usersRepo:     r,
	}
}

func (u *usersUseCase) Register(ctx context.Context, req *usecase.RegisterRequest) (*usecase.RegistrationResult, error) {
	existingUser, err := u.usersRepo.GetUserByName(ctx, req.Username)
	if existingUser != nil {
		return nil, usecase.ErrUsernameTaken{Username: existingUser.Username}
	}

	saltBytes := make([]byte, SaltForKeysSize)
	rand.Seed(time.Now().UnixNano())
	rand.Read(saltBytes)

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &usecase.FullUser{
		Username:                      req.Username,
		PasswordHash:                  passwordHashBytes,
		SaltForKeys:                   saltBytes,
		PublicKeyForSharing:           req.PublicKeyForSharing,
		EncryptedPrivateKeyForSharing: req.EncryptedPrivateKeyForSharing,
	}

	user, err = u.usersRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	tokenStr, err := u.tokensService.GenerateToken(user.Id)
	registrationResult := &usecase.RegistrationResult{
		Token:       tokenStr,
		SaltForKeys: string(saltBytes),
	}

	return registrationResult, nil
}

func (u *usersUseCase) Login(ctx context.Context, request *usecase.LoginRequest) (*usecase.AuthResult, error) {
	user, err := u.usersRepo.GetUserByName(ctx, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrUserNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password))
	if err != nil {
		return nil, err
	}

	token, err := u.tokensService.GenerateToken(user.Id)
	if err != nil {
		return nil, err
	}

	authResult := &usecase.AuthResult{Token: token}
	return authResult, nil
}

func (u *usersUseCase) RefreshToken(ctx context.Context, token string) (*usecase.AuthResult, error) {
	newToken, err := u.tokensService.RefreshToken(token)
	if err != nil {
		return nil, err
	}

	return &usecase.AuthResult{Token: newToken}, nil
}

func (u *usersUseCase) GetFullUser(ctx context.Context, userId uuid.UUID) (*usecase.FullUser, error) {
	user, err := u.usersRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersUseCase) GetUserById(ctx context.Context, userId uuid.UUID) (*usecase.ShortUser, error) {
	user, err := u.GetFullUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	shortUser := &usecase.ShortUser{
		Id:                  user.Id,
		Username:            user.Username,
		PublicKeyForSharing: user.PublicKeyForSharing,
	}
	return shortUser, nil
}

func (u *usersUseCase) GetUserByName(ctx context.Context, username string) (*usecase.ShortUser, error) {
	fullUser, err := u.usersRepo.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}

	user := &usecase.ShortUser{
		Id:                  fullUser.Id,
		Username:            fullUser.Username,
		PublicKeyForSharing: fullUser.PublicKeyForSharing,
	}

	return user, nil
}
