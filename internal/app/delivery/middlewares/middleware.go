package middlewares

import (
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
)

type MiddleWare struct {
	authSvc  service.AuthService
	tokensvc token.TokenService
}

func NewMiddelwWare(authSvc service.AuthService, tokensvc token.TokenService) *MiddleWare {
	return &MiddleWare{
		authSvc:  authSvc,
		tokensvc: tokensvc,
	}
}
