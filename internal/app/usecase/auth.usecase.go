package usecase

import (
	"errors"
	"strconv"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
	"onthemat/pkg/auth/store"
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

	SendEmailResetPassword(ctx *fasthttp.RequestCtx, email string) error
	CheckDuplicatedEmail(ctx *fasthttp.RequestCtx, email string) error
	VerifiedEmail(ctx *fasthttp.RequestCtx, email string, authKey string) error
	RefreshToken(ctx *fasthttp.RequestCtx, refreshToken string) error
}

type authUseCase struct {
	tokenSvc token.TokenService
	authSvc  service.AuthService
	userRepo repository.UserRepository
	store    store.Store
	config   *config.Config
}

func NewAuthUseCase(
	tokenSvc token.TokenService,
	userRepo repository.UserRepository,
	authsvc service.AuthService,
	store store.Store,
	config *config.Config,
) AuthUseCase {
	return &authUseCase{
		tokenSvc: tokenSvc,
		authSvc:  authsvc,
		userRepo: userRepo,
		store:    store,
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
	user, err := a.userRepo.GetByEmailPassword(ctx, &ent.User{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("이메일 혹은 비밀번호를 다시 확인해주세요. ")
		}
		panic(err)
	}

	if !user.IsEmailVerified {
		return nil, errors.New("이메일 인증이 필요합니다. ")
	}

	userType := ""

	if user.Type != nil {
		userType = string(*user.Type)
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, _ := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.RefreshTokenExpired)
	access, _ := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.AccessTokenExpired)

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

		kakaoId := strconv.FormatUint(uint64(kakaoInfo.Id), 10)

		user.SocialKey = &kakaoId

	} else if socialName == "google" {

		googleInfo, err := a.authSvc.GetGoogleInfo(code)
		if err != nil {
			panic(err)
		}

		googleId := strconv.FormatUint(uint64(googleInfo.Sub), 10)
		user.SocialKey = &googleId
		user.Email = googleInfo.Email
	} else if socialName == "naver" {
		naverInfo, err := a.authSvc.GetNaverInfo(code)
		if err != nil {
			panic(err)
		}

		user.SocialKey = &naverInfo.Id
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
	isExist, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		panic(err)
	}

	if isExist {
		return errors.New("이미 존재하는 이메일입니다. ")
	}

	_, err = a.userRepo.Create(ctx, &ent.User{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		panic(err)
	}

	key := a.authSvc.GenerateRandomString()
	if err := a.store.Set(ctx, body.Email, key, time.Duration(time.Hour*24)); err != nil {
		panic(err)
	}

	a.authSvc.SendEmailVerifiedUser(body.Email, key)

	return err
}

func (a *authUseCase) CheckDuplicatedEmail(ctx *fasthttp.RequestCtx, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		panic(err)
	}

	if isExist {
		return errors.New("이미 존재하는 이메일입니다. ")
	}

	return nil
}

func (a *authUseCase) VerifiedEmail(ctx *fasthttp.RequestCtx, email string, authKey string) error {
	key := a.store.Get(ctx, email)

	if key != authKey {
		return errors.New("올바르지 않은 인증키 입니다. ")
	}

	if err := a.userRepo.UpdateEmailVerified(ctx, email); err != nil {
		panic(err)
	}

	a.store.Set(ctx, email, "", 0)

	return nil
}

func (a authUseCase) SendEmailResetPassword(ctx *fasthttp.RequestCtx, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		panic(err)
	}

	if !isExist {
		return errors.New("존재하지 않는 이메일입니다. ")
	}

	u := &ent.User{
		Email:        email,
		TempPassword: a.authSvc.GenerateRandomPassword(),
	}
	a.userRepo.UpdateTempPassword(ctx, u)
	a.authSvc.SendEmailResetPassword(u)

	return nil
	// 이메일로 패스워드 찾아서,
}
