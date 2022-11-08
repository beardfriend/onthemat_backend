package usecase

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"onthemat/internal/app/common"
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
	Refresh(ctx *fasthttp.RequestCtx, authorizationHeader []byte) (*RefreshResult, error)
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
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiredAt  time.Time `json:"accessTokenExpiredAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiredAt time.Time `json:"refreshTokenExpiredAt"`
}

func (a *authUseCase) Login(ctx *fasthttp.RequestCtx, body *transport.LoginBody) (*LoginResult, error) {
	defer ctx.Done()

	sha := sha256.New()
	sha.Write([]byte(a.config.Secret.Password))
	sha.Write([]byte(body.Password))
	hashPassword := fmt.Sprintf("%x", sha.Sum(nil))

	user, err := a.userRepo.GetByEmailPassword(ctx, &ent.User{
		Email:    body.Email,
		Password: hashPassword,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewBadRequestError("이메일 혹은 비밀번호를 다시 확인해주세요.")
		}
		return nil, err
	}

	if !user.IsEmailVerified {
		return nil, common.NewBadRequestError("이메일 인증이 필요합니다.")
	}

	userType := ""

	if user.Type != nil {
		userType = string(*user.Type)
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, err := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.RefreshTokenExpired)
	if err != nil {
		return nil, err
	}

	if err := a.store.Set(ctx, uid, strconv.Itoa(user.ID), time.Duration(a.config.JWT.RefreshTokenExpired)*time.Minute); err != nil {
		return nil, err
	}

	access, err := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:           access,
		AccessTokenExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
		RefreshToken:          refresh,
		RefreshTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.RefreshTokenExpired),
	}, nil
}

func (a *authUseCase) SocialLogin(ctx *fasthttp.RequestCtx, socialName, code string) (*LoginResult, error) {
	user := new(ent.User)
	// TODO: 일반 에러로 빼기.
	if socialName != "kakao" && socialName != "google" && socialName != "naver" {
		return nil, common.NewBadRequestError("올바르지 않은 소셜 이름입니다.")
	}

	if socialName == "kakao" {

		kakaoInfo, err := a.authSvc.GetKakaoInfo(code)
		if err != nil {
			return nil, err
		}

		kakaoId := strconv.FormatUint(uint64(kakaoInfo.Id), 10)

		user.SocialKey = &kakaoId

	} else if socialName == "google" {

		googleInfo, err := a.authSvc.GetGoogleInfo(code)
		if err != nil {
			return nil, err
		}

		googleId := strconv.FormatUint(uint64(googleInfo.Sub), 10)
		user.SocialKey = &googleId
		user.Email = googleInfo.Email
	} else if socialName == "naver" {
		naverInfo, err := a.authSvc.GetNaverInfo(code)
		if err != nil {
			return nil, err
		}

		user.SocialKey = &naverInfo.Id
		user.Email = naverInfo.Email
	}

	user.SocialName = &socialName

	checkedUser, err := a.userRepo.GetBySocialKey(ctx, user)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	// 유저가 없으면 회원 정보 생성
	if checkedUser == nil {
		user, err = a.userRepo.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	userType := ""

	if user.Type != nil {
		userType = string(*user.Type)
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, err := a.tokenSvc.GenerateToken(uid, user.ID, socialName, userType, a.config.JWT.RefreshTokenExpired)
	if err != nil {
		return nil, err
	}

	if err := a.store.Set(ctx, uid, strconv.Itoa(user.ID), time.Duration(a.config.JWT.RefreshTokenExpired)*time.Minute); err != nil {
		return nil, err
	}

	access, err := a.tokenSvc.GenerateToken(uid, user.ID, socialName, userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:           access,
		AccessTokenExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
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
		return err
	}

	if isExist {
		return common.NewConflictError("이미 존재하는 이메일입니다.")
	}

	termAgreeAt := time.Now()

	sha := sha256.New()
	sha.Write([]byte(a.config.Secret.Password))
	sha.Write([]byte(body.Password))
	hashPassword := fmt.Sprintf("%x", sha.Sum(nil))

	_, err = a.userRepo.Create(ctx, &ent.User{
		Email:       body.Email,
		Password:    hashPassword,
		Nickname:    &body.NickName,
		TermAgreeAt: &termAgreeAt,
	})

	if err != nil {
		return err
	}

	key := a.authSvc.GenerateRandomString()
	if err := a.store.Set(ctx, body.Email, key, time.Duration(time.Hour*24)); err != nil {
		return err
	}

	go a.authSvc.SendEmailVerifiedUser(body.Email, key)

	return err
}

func (a *authUseCase) CheckDuplicatedEmail(ctx *fasthttp.RequestCtx, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if isExist {
		return common.NewConflictError("이미 존재하는 이메일입니다.")
	}

	return nil
}

func (a *authUseCase) VerifiedEmail(ctx *fasthttp.RequestCtx, email string, authKey string) error {
	key := a.store.Get(ctx, email)

	if key != authKey {
		return common.NewBadRequestError("올바르지 않은 인증키입니다.")
	}

	u, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if u.IsEmailVerified {
		return common.NewConflictError("이미 인증된 유저입니다.")
	}

	if err := a.userRepo.UpdateEmailVerifeid(ctx, u.ID); err != nil {
		return err
	}

	if err := a.store.Del(ctx, email); err != nil {
		return err
	}

	return nil
}

func (a *authUseCase) SendEmailResetPassword(ctx *fasthttp.RequestCtx, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if !isExist {
		return common.NewBadRequestError("존재하지 않는 이메일입니다.")
	}

	sha := sha256.New()
	sha.Write([]byte(a.config.Secret.Password))
	sha.Write([]byte(a.authSvc.GenerateRandomPassword()))
	hashPassword := fmt.Sprintf("%x", sha.Sum(nil))

	u := &ent.User{
		Email:        email,
		TempPassword: hashPassword,
	}
	if err := a.userRepo.UpdateTempPassword(ctx, u); err != nil {
		return err
	}
	if err := a.authSvc.SendEmailResetPassword(u); err != nil {
		return err
	}

	return nil
	// 이메일로 패스워드 찾아서,
}

type RefreshResult struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiredAt time.Time `json:"accessTokenExpiredAt"`
}

func (a *authUseCase) Refresh(ctx *fasthttp.RequestCtx, authorizationHeader []byte) (*RefreshResult, error) {
	refreshToken, err := a.authSvc.ExtractTokenFromHeader(string(authorizationHeader))
	if err != nil {
		return nil, common.NewBadRequestError("헤더를 확인해주세요.")
	}

	claim := new(token.TokenClaim)
	if err := a.tokenSvc.ParseToken(refreshToken, claim); err != nil {
		return nil, err
	}

	val := a.store.Get(ctx, claim.Uuid)
	if val != strconv.Itoa(claim.UserId) {
		return nil, common.NewAuthenticationFailedError("잘못된 토큰입니다.")
	}

	u, err := a.userRepo.Get(ctx, claim.UserId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewNotFoundError("존재하지 않는 유저입니다.")
		}
		return nil, err
	}

	uid := uuid.New().String()

	userType := ""
	if u.Type != nil {
		userType = string(*u.Type)
	}

	loginType := ""
	if u.SocialName != nil || u.SocialName != &loginType {
		loginType = "social"
	}

	access, err := a.tokenSvc.GenerateToken(uid, u.ID, loginType, userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:          access,
		AccessTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
	}, nil
}
