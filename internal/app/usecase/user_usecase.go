package usecase

import (
	"context"
	"errors"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"

	"github.com/lib/pq"
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
			// unwrap
			err = err.(*ent.ConstraintError).Unwrap()
			err = errors.Unwrap(err)

			code := err.(*pq.Error).Code
			if code == "23503" {
				err = ex.NewConflictError(ex.ErrYogaAlreadyRegisted, nil)
				return
			}
			if code == "23505" {
				err = ex.NewConflictError(ex.ErrYogaDoseNotExist, nil)
				return
			}
			err = ex.NewBadRequestError(ex.ErrYogaIdsInvliad, nil)
			return err
		}
		return
	}

	return
}
