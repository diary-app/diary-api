package usecase

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, request *RegisterRequest) (*RegistrationResult, error)
	Login(ctx context.Context, request *LoginRequest) (*AuthResult, error)
	RefreshToken(ctx context.Context, token string) (*AuthResult, error)
}
