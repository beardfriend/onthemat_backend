package usecase

import (
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type AuthUseCase interface {
	SignUp(ctx *fasthttp.RequestCtx, body *transport.SignUpBody) error
	SocialSignUp(ctx *fasthttp.RequestCtx, body *transport.SocialSignUpBody) error
	KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string
	KakaoLogin(ctx *fasthttp.RequestCtx, code string) (string, string)
}

type authUseCase struct {
	tokenSvc token.TokenService
	authSvc  service.AuthService
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthUseCase(
	tokenSvc token.TokenService,
	userRepo repository.UserRepository,
	authsvc service.AuthService,
	config *config.Config,
) AuthUseCase {
	return &authUseCase{
		tokenSvc: tokenSvc,
		authSvc:  authsvc,
		userRepo: userRepo,
		config:   config,
	}
}

func (a *authUseCase) KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string {
	return a.authSvc.GetKakaoRedirectUrl()
}

func (a *authUseCase) KakaoLogin(ctx *fasthttp.RequestCtx, code string) (string, string) {
	// get kakao Info FROM kakao
	kakaoInfo, err := a.authSvc.GetKakaoInfo(code)
	if err != nil {
		panic(err)
	}

	u := &ent.User{
		SocialKey:  int(kakaoInfo.Id),
		SocialName: "kakao",
	}
	user, err := a.userRepo.GetBySocialKey(ctx, u)
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}

	// 유저가 없으면 회원 정보 생성
	if user == nil {
		user, err = a.userRepo.Create(ctx, u)
		if err != nil {
			panic(err)
		}
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, _ := a.tokenSvc.GenerateToken(uid, user.ID, "kakao", "", a.config.JWT.RefreshTokenExpired)
	access, _ := a.tokenSvc.GenerateToken(uid, user.ID, "kakao", "", a.config.JWT.AccessTokenExpired)
	return access, refresh
}

func (a *authUseCase) SocialSignUp(ctx *fasthttp.RequestCtx, body *transport.SocialSignUpBody) error {
	var termAgreeAt *time.Time
	if body.TermAgree {
		*termAgreeAt = time.Now()
	}

	_, err := a.userRepo.Update(ctx, &ent.User{
		Email:       body.Email,
		Nickname:    body.NickName,
		TermAgreeAt: *termAgreeAt,
	})
	return err
}

func (a *authUseCase) SignUp(ctx *fasthttp.RequestCtx, body *transport.SignUpBody) error {
	_, err := a.userRepo.Create(ctx, &ent.User{
		Email:    body.Email,
		Password: body.Password,
	})
	return err
}
