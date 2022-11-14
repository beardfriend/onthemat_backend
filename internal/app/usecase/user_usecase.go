package usecase

import (
	"context"

	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"
)

type UserUseCase interface {
	GetMe(ctx context.Context, id int) (*ent.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GetMe(ctx context.Context, id int) (*ent.User, error) {
	return u.userRepo.Get(ctx, id)
}
