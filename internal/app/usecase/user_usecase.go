package usecase

import (
	"context"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/transport/request"
	"onthemat/pkg/ent"
)

type UserUseCase interface {
	GetMe(ctx context.Context, id int) (result *ent.User, err error)
	Update(ctx context.Context, reqBody *request.UserUpdateBody, id int) (err error)
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

func (u *userUseCase) Update(ctx context.Context, reqBody *request.UserUpdateBody, id int) (err error) {
	_, err = u.userRepo.Update(ctx, &ent.User{
		ID:       id,
		Nickname: &reqBody.Nickname,
		PhoneNum: &reqBody.PhoneNum,
	})

	if err != nil {
		return
	}
	return
}
