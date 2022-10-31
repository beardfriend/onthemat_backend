package usecase

import (
	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"

	"github.com/valyala/fasthttp"
)

type UserUseCase interface {
	GetMe(ctx *fasthttp.RequestCtx, id int) (*ent.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) GetMe(ctx *fasthttp.RequestCtx, id int) (*ent.User, error) {
	return u.userRepo.Get(ctx, id)
}
