package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	"onthemat/internal/app/config"
	"onthemat/internal/app/dto"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service/token"
	"onthemat/pkg/ent"
	"onthemat/pkg/kakao"

	"github.com/valyala/fasthttp"
)

type AuthUseCase interface {
	SignUp(ctx *fasthttp.RequestCtx, body *dto.SignUpBody) error
	KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string
	KakaoLogin(ctx *fasthttp.RequestCtx, code string) error
}

type authUseCase struct {
	tokenSvc token.TokenService
	kakao    *kakao.Kakao
	userRepo repository.UserRepository
}

func NewAuthUseCase(config *config.Config, tokenSvc token.TokenService, userRepo repository.UserRepository) AuthUseCase {
	kakao := kakao.NewKakao(config)
	return &authUseCase{
		tokenSvc: tokenSvc,
		userRepo: userRepo,
		kakao:    kakao,
	}
}

func (a *authUseCase) KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string {
	resp := a.kakao.Authorize()
	r := resp.Header.Peek("Location")
	return string(r)
}

func (a *authUseCase) KakaoLogin(ctx *fasthttp.RequestCtx, code string) error {
	resp := a.kakao.GetToken(code)

	if resp.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(resp.Body(), body)

		return errors.New(body.Error + body.ErrorCode)
	}

	body := new(kakao.GetTokenSuccessBody)
	json.Unmarshal(resp.Body(), body)

	respInfo := a.kakao.GetUserInfo(body.AccessToken)

	if respInfo.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(resp.Body(), body)

		return errors.New(body.Error + body.ErrorCode)
	}

	userInfoBody := new(kakao.GetUserInfoSuccessBody)
	json.Unmarshal(respInfo.Body(), respInfo)
	key := fmt.Sprintf("%v", userInfoBody.Id)

	u := &ent.User{SocialKey: key, SocialName: "kakao"}
	exist, err := a.userRepo.FindBySocialKey(ctx, u)
	if err != nil {
		return err
	}

	if !exist {
		if err := a.userRepo.Create(ctx, u); err != nil {
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
