package users

import (
	"context"
	"diary-api/internal/usecase"
	"github.com/google/uuid"
)

type usersUseCase struct {
	usersRepo   usecase.UsersRepository
	diariesRepo usecase.DiaryRepository
}

func NewUseCase(r usecase.UsersRepository) usecase.UsersUseCase {
	return &usersUseCase{
		usersRepo: r,
	}
}

func (u *usersUseCase) GetFullUser(ctx context.Context, userID uuid.UUID) (*usecase.FullUser, error) {
	user, err := u.usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (*usecase.ShortUser, error) {
	user, err := u.GetFullUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	shortUser := &usecase.ShortUser{
		ID:                  user.ID,
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
		ID:                  fullUser.ID,
		Username:            fullUser.Username,
		PublicKeyForSharing: fullUser.PublicKeyForSharing,
	}

	return user, nil
}
