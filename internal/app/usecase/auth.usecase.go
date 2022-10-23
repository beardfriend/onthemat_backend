package usecase

import (
	"onthemat/internal/app/config"
	"onthemat/internal/app/dto"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service/token"
)

type AuthUseCase interface {
	Extract()
}

type authUseCase struct {
	tokenSvc token.TokenService
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthUseCase() {
}

func (a *authUseCase) SignUp(query dto.LoginRequestQuery) {
}

func Logout() {
}
