package usecases

import (
	"context"
	"main/src/features/user/domain/entities"
	repositories "main/src/features/user/domain/repositories"
)

type UserUsecase interface {
	GetUserById(ctx context.Context, id int64) (*entities.User, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
}

func (u *userUsecase) GetUserById(ctx context.Context, id int64) (*entities.User, error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserUsecase(userRepo repositories.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}
