package usecase

import (
	"context"

	"onthemat/internal/app/dto"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service/token"
	"onthemat/pkg/ent"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, body *dto.SignUpBody) error
}

type authUseCase struct {
	tokenSvc token.TokenService
	userRepo repository.UserRepository
}

func NewAuthUseCase(tokenSvc token.TokenService, userRepo repository.UserRepository) AuthUseCase {
	return authUseCase{
		tokenSvc: tokenSvc,
		userRepo: userRepo,
	}
}

func (a *authUseCase) SignUp(ctx context.Context, body *dto.SignUpBody) error {
	data := &ent.User{
		Email:    body.Email,
		Password: body.Password,
	}
	return a.userRepo.Create(ctx, data)
}
