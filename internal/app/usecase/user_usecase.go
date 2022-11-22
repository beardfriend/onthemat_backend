package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"
)

type UserUseCase interface {
	GetMe(ctx context.Context, id int) (*ent.User, error)
	AddYoga(ctx context.Context, id int, yogaIds []int) (err error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GetMe(ctx context.Context, id int) (result *ent.User, err error) {
	result, err = u.userRepo.Get(ctx, id)
	if ent.IsNotFound(err) {
		err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
		return
	}
	return
}

func (u *userUseCase) AddYoga(ctx context.Context, id int, yogaIds []int) (err error) {
	err = u.userRepo.AddYoga(ctx, id, yogaIds)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewBadRequestError(ex.ErrYogaIdsInvliad, nil)
			return
		}
		return
	}

	return
}
