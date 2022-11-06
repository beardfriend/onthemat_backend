package usecase

import (
	"crypto/sha256"
	"errors"
	"strconv"
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
	Login(ctx *fasthttp.RequestCtx, body *transport.LoginBody) (*LoginResult, error)
	SocialSignUp(ctx *fasthttp.RequestCtx, body *transport.SocialSignUpBody) error
	SocialLogin(ctx *fasthttp.RequestCtx, socialName, code string) (*LoginResult, error)
	KakaoRedirectUrl(ctx *fasthttp.RequestCtx) string
	NaverRedirectUrl(ctx *fasthttp.RequestCtx) string
	GoogleRedirectUrl(ctx *fasthttp.RequestCtx) string
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

func (a *authUseCase) NaverRedirectUrl(ctx *fasthttp.RequestCtx) string {
	return a.authSvc.GetNaverRedirectUrl()
}

type LoginResult struct {
	AccessToken           string
	AccessToeknExpiredAt  time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
}

func (a *authUseCase) Login(ctx *fasthttp.RequestCtx, body *transport.LoginBody) (*LoginResult, error) {
	body.Password = string(sha256.New().Sum([]byte(body.Password)))
	user, err := a.userRepo.GetByEmailPassword(ctx, &ent.User{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		panic(err)
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, _ := a.tokenSvc.GenerateToken(uid, user.ID, "normal", string(*user.Type), a.config.JWT.RefreshTokenExpired)
	access, _ := a.tokenSvc.GenerateToken(uid, user.ID, "normal", string(*user.Type), a.config.JWT.AccessTokenExpired)

	return &LoginResult{
		AccessToken:           access,
		AccessToeknExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
		RefreshToken:          refresh,
		RefreshTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.RefreshTokenExpired),
	}, nil
}

func (a *authUseCase) SocialLogin(ctx *fasthttp.RequestCtx, socialName, code string) (*LoginResult, error) {
	user := new(ent.User)
	if socialName != "kakao" && socialName != "google" && socialName != "naver" {
		return nil, errors.New("올바르지 않은 socialName입니다. ")
	}

	if socialName == "kakao" {

		kakaoInfo, err := a.authSvc.GetKakaoInfo(code)
		if err != nil {
			panic(err)
		}

		kakaoId := int(kakaoInfo.Id)

		user.SocialKey = &kakaoId

	} else if socialName == "google" {

		googleInfo, err := a.authSvc.GetGoogleInfo(code)
		if err != nil {
			panic(err)
		}

		googleId := int(googleInfo.Sub)
		user.SocialKey = &googleId
		user.Email = googleInfo.Email
	} else if socialName == "naver" {
		naverInfo, err := a.authSvc.GetNaverInfo(code)
		if err != nil {
			panic(err)
		}

		naverId, _ := strconv.Atoi(naverInfo.Id)
		user.SocialKey = string(&naverId
		user.Email = naverInfo.Email
	}
	user.SocialName = &socialName

	checkedUser, err := a.userRepo.GetBySocialKey(ctx, user)
	if err != nil && !ent.IsNotFound(err) {
		panic(err)
	}

	// 유저가 없으면 회원 정보 생성
	if checkedUser == nil {
		user, err = a.userRepo.Create(ctx, user)
		if err != nil {
			panic(err)
		}
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, _ := a.tokenSvc.GenerateToken(uid, user.ID, socialName, string(*user.Type), a.config.JWT.RefreshTokenExpired)
	access, _ := a.tokenSvc.GenerateToken(uid, user.ID, socialName, string(*user.Type), a.config.JWT.AccessTokenExpired)

	return &LoginResult{
		AccessToken:           access,
		AccessToeknExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
		RefreshToken:          refresh,
		RefreshTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.RefreshTokenExpired),
	}, nil
}

func (a *authUseCase) GoogleRedirectUrl(ctx *fasthttp.RequestCtx) string {
	return a.authSvc.GetGoogleRedirectUrl()
}

func (a *authUseCase) SocialSignUp(ctx *fasthttp.RequestCtx, body *transport.SocialSignUpBody) error {
	termAgreeAt := time.Now()
	_, err := a.userRepo.Update(ctx, &ent.User{
		ID:          body.UserID,
		Email:       body.Email,
		Nickname:    &body.NickName,
		TermAgreeAt: &termAgreeAt,
		Type:        nil,
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
