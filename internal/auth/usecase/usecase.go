package usecase

import (
	"context"
	"diary-api/internal/usecase"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	storage      storage
	tokenService tokenService
}

func New(usersRepo storage, tokenService tokenService) *UseCase {
	return &UseCase{storage: usersRepo, tokenService: tokenService}
}

func (u *UseCase) Register(ctx context.Context, req *usecase.RegisterRequest) (*usecase.RegistrationResult, error) {
	existingUser, err := u.storage.GetUserByName(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, usecase.ErrUsernameTaken{Username: existingUser.Username}
	}
	if err != nil {
		return nil, err
	}

	user, diary, err := mapRegisterRequestToUserAndDiary(req)
	if err != nil {
		return nil, err
	}

	user, diary, err = u.storage.CreateUser(ctx, user, diary)
	if err != nil {
		return nil, err
	}

	tokenStr, err := u.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	registrationResult := &usecase.RegistrationResult{
		Token:   tokenStr,
		DiaryID: diary.ID,
	}
	return registrationResult, nil
}

func mapRegisterRequestToUserAndDiary(req *usecase.RegisterRequest) (*usecase.FullUser, *usecase.Diary, error) {
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}
	user := &usecase.FullUser{
		Username:                      req.Username,
		PasswordHash:                  passwordHashBytes,
		SaltForKeys:                   []byte(req.MasterKeySalt),
		PublicKeyForSharing:           req.PublicKeyForSharing,
		EncryptedPrivateKeyForSharing: req.EncryptedPrivateKeyForSharing,
	}
	keys := make([]usecase.DiaryKey, 1)
	keys[0] = usecase.DiaryKey{
		EncryptedKey: req.EncryptedDiaryKey,
	}
	diary := &usecase.Diary{
		Name: fmt.Sprintf("Дневник пользователя %s", req.Username),
		Keys: keys,
	}
	return user, diary, nil
}

func (u *UseCase) Login(ctx context.Context, request *usecase.LoginRequest) (*usecase.AuthResult, error) {
	user, err := u.storage.GetUserByName(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, usecase.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password))
	if err != nil {
		return nil, err
	}

	token, err := u.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	authResult := &usecase.AuthResult{Token: token}
	return authResult, nil
}

func (u *UseCase) RefreshToken(_ context.Context, token string) (*usecase.AuthResult, error) {
	newToken, err := u.tokenService.RefreshToken(token)
	if err != nil {
		return nil, err
	}

	return &usecase.AuthResult{Token: newToken}, nil
}
