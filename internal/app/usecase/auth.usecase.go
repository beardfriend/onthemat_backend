package usecase

import (
	"onthemat/internal/app/dto"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/pkg/ent"

	"github.com/valyala/fasthttp"
)

type AuthUseCase interface {
	SignUp(ctx *fasthttp.RequestCtx, body *dto.SignUpBody) error
	KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string
	KakaoLogin(ctx *fasthttp.RequestCtx, code string) error
}

type authUseCase struct {
	tokenSvc token.TokenService
	authSvc  service.AuthService
	userRepo repository.UserRepository
}

func NewAuthUseCase(tokenSvc token.TokenService, userRepo repository.UserRepository, authsvc service.AuthService) AuthUseCase {
	return &authUseCase{
		tokenSvc: tokenSvc,
		authSvc:  authsvc,
		userRepo: userRepo,
	}
}

func (a *authUseCase) KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string {
	return a.authSvc.GetKakaoRedirectUrl()
}

func (a *authUseCase) KakaoLogin(ctx *fasthttp.RequestCtx, code string) error {
	kakaoId, err := a.authSvc.GetKakaoID(code)
	if err != nil {
		return err
	}

	u := &ent.User{SocialKey: kakaoId, SocialName: "kakao"}
	exist, err := a.userRepo.FindBySocialKey(ctx, u)
	if err != nil {
		return err
	}

	if !exist {
		err = a.userRepo.Create(ctx, u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *authUseCase) SignUp(ctx *fasthttp.RequestCtx, body *dto.SignUpBody) error {
	user := &ent.User{
		Email:    body.Email,
		Password: body.Password,
	}
	return a.userRepo.Create(ctx, user)
}
