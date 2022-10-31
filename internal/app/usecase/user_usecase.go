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

func NewUserUseCase(
	userRepo repository.UserRepository,
) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

func (a *authUseCase) GetMe(ctx *fasthttp.RequestCtx, id int) (*ent.User, error) {
	return a.userRepo.Get(ctx, id)
}
