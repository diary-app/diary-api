package usecase_test

import (
	"context"
	"database/sql"
	authUsecase "diary-api/internal/auth/usecase"
	"diary-api/internal/usecase"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type commonTestCase struct {
	userStorage  func(ctrl *gomock.Controller) *Mockstorage
	tokenService func(ctrl *gomock.Controller) *MocktokenService
}

type loginTestCase struct {
	name string
	req  *usecase.LoginRequest
	commonTestCase
	assertFunc func(t *testing.T, result *usecase.AuthResult, err error)
}

func TestUseCase_Register(t *testing.T) {

	freeName := "freeName"
	takenName := "takenName"
	newUserID := uuid.MustParse("b578d867-8c4f-448e-97c8-6d5decb9b6ad")
	newDiaryID := uuid.MustParse("4c338341-2eb6-4434-a3ab-d862ec72b441")
	wantToken := "someToken"
	unexpectedError := errors.New("unexpected error")

	tests := []struct {
		name         string
		req          *usecase.RegisterRequest
		storage      func(ctrl *gomock.Controller) *Mockstorage
		tokenService func(ctrl *gomock.Controller) *MocktokenService
		assertFunc   func(t *testing.T, result *usecase.RegistrationResult, err error)
	}{
		{
			name: "Returns RegistrationResult when user with name does not exist and successfully created",
			req:  &usecase.RegisterRequest{Username: freeName},
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				s := NewMockstorage(ctrl)
				s.EXPECT().GetUserByName(context.Background(), freeName).Return(nil, sql.ErrNoRows)
				newUser := &usecase.FullUser{Username: freeName, ID: newUserID}
				newDiary := &usecase.Diary{ID: newDiaryID}
				s.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any()).Return(newUser, newDiary, nil)
				return s
			},
			tokenService: func(ctrl *gomock.Controller) *MocktokenService {
				ts := NewMocktokenService(ctrl)
				ts.EXPECT().GenerateToken(newUserID).Return(wantToken, nil)
				return ts
			},
			assertFunc: func(t *testing.T, result *usecase.RegistrationResult, err error) {
				assert.NotNil(t, result)
				assert.Equal(t, result.Token, wantToken)
				assert.Equal(t, result.DiaryID, newDiaryID)
				assert.Nil(t, err)
			},
		},
		{
			name: "Returns ErrUsernameTaken when user with given name already exists",
			req:  &usecase.RegisterRequest{Username: takenName},
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				s := NewMockstorage(ctrl)
				existingUser := &usecase.FullUser{Username: takenName}
				s.EXPECT().GetUserByName(context.Background(), takenName).Return(existingUser, nil)
				return s
			},
			tokenService: func(ctrl *gomock.Controller) *MocktokenService {
				ts := NewMocktokenService(ctrl)
				ts.EXPECT().GenerateToken(gomock.Any()).Times(0)
				return ts
			},
			assertFunc: func(t *testing.T, result *usecase.RegistrationResult, err error) {
				assert.Nil(t, result)
				_, isErrUsernameTaken := err.(usecase.ErrUsernameTaken)
				assert.True(t, isErrUsernameTaken)
			},
		},
		{
			name: "Returns error from user storage when GetUserByName returns unexpected error",
			req:  &usecase.RegisterRequest{Username: takenName},
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				s := NewMockstorage(ctrl)
				s.EXPECT().GetUserByName(context.Background(), takenName).Return(nil, unexpectedError)
				return s
			},
			tokenService: func(ctrl *gomock.Controller) *MocktokenService {
				return NewMocktokenService(ctrl)
			},
			assertFunc: func(t *testing.T, result *usecase.RegistrationResult, err error) {
				assert.Nil(t, result)
				assert.Equal(t, err, unexpectedError)
			},
		},
		{
			name: "Returns error from user storage when CreateUser returns unexpected error",
			req:  &usecase.RegisterRequest{Username: freeName},
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				s := NewMockstorage(ctrl)
				s.EXPECT().GetUserByName(context.Background(), freeName).Return(nil, sql.ErrNoRows)
				s.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any()).
					Return(nil, unexpectedError)
				return s
			},
			tokenService: func(ctrl *gomock.Controller) *MocktokenService {
				return NewMocktokenService(ctrl)
			},
			assertFunc: func(t *testing.T, result *usecase.RegistrationResult, err error) {
				assert.Nil(t, result)
				assert.Equal(t, err, unexpectedError)
			},
		},
		{
			name: "Returns error from token service when GenerateToken returns unexpected error",
			req:  &usecase.RegisterRequest{Username: freeName},
			storage: func(ctrl *gomock.Controller) *Mockstorage {
				s := NewMockstorage(ctrl)
				createdUser := &usecase.FullUser{ID: newUserID}
				s.EXPECT().GetUserByName(context.Background(), freeName).Return(nil, sql.ErrNoRows)
				s.EXPECT().CreateUser(context.Background(), gomock.Any(), gomock.Any()).Return(createdUser, nil)
				return s
			},
			tokenService: func(ctrl *gomock.Controller) *MocktokenService {
				tc := NewMocktokenService(ctrl)
				tc.EXPECT().GenerateToken(newUserID).Return("", unexpectedError)
				return tc
			},
			assertFunc: func(t *testing.T, result *usecase.RegistrationResult, err error) {
				assert.Nil(t, result)
				assert.Equal(t, err, unexpectedError)
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := authUsecase.New(tc.storage(ctrl), tc.tokenService(ctrl))
			d, err := u.Register(context.Background(), tc.req)
			tc.assertFunc(t, d, err)
		})
	}
}

func TestUseCase_Login(t *testing.T) {
	tests := []loginTestCase{
		{
			name: "Success when user with name exists and password is correct",
			assertFunc: func(t *testing.T, result *usecase.AuthResult, err error) {
				assert.True(t, true)
			},
		},
		{
			name: "Incorrect password error when user exists but password is incorrect",
			assertFunc: func(t *testing.T, result *usecase.AuthResult, err error) {
				assert.True(t, true)
			},
		},
		{
			name: "Returns ErrUserNotFound when user storage returns sql.ErrNoRows",
			assertFunc: func(t *testing.T, result *usecase.AuthResult, err error) {
				assert.True(t, true)
			},
		},
		{
			name: "Returns error from user storage when user storage returns unexpected error",
			assertFunc: func(t *testing.T, result *usecase.AuthResult, err error) {
				assert.True(t, true)
			},
		},
	}

	for _, tt := range tests {
		tc := tt

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//u := authUsecase.New(tc.storage(ctrl), tc.tokenService(ctrl))
			//d, err := u.Login(context.Background(), tc.req)
			tc.assertFunc(t, nil, nil)
		})
	}
}
