package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"onthemat/internal/app/common"
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/model"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/auth/store"
	"onthemat/pkg/ent"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, body *transport.SignUpBody) error
	Login(ctx context.Context, body *transport.LoginBody) (*LoginResult, error)
	SocialSignUp(ctx context.Context, body *transport.SocialSignUpBody) error
	SocialLogin(ctx context.Context, socialName model.SocialType, code string) (*LoginResult, error)
	KakaoRedirectUrl(ctx context.Context) string
	NaverRedirectUrl(ctx context.Context) string
	GoogleRedirectUrl(ctx context.Context) string

	// LogOut
	// PasswordChange
	// 탈퇴

	// 입력받은 이메일로 임시 비밀번호 전송하는 모듈
	SendEmailResetPassword(ctx context.Context, email string) error
	CheckDuplicatedEmail(ctx context.Context, email string) error
	VerifiyEmail(ctx context.Context, email string, authKey string, issuedAt string) error
	Refresh(ctx context.Context, authorizationHeader []byte) (*RefreshResult, error)
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

func (a *authUseCase) KakaoRedirectUrl(ctx context.Context) string {
	return a.authSvc.GetKakaoRedirectUrl()
}

func (a *authUseCase) NaverRedirectUrl(ctx context.Context) string {
	return a.authSvc.GetNaverRedirectUrl()
}

type LoginResult struct {
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiredAt  time.Time `json:"accessTokenExpiredAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiredAt time.Time `json:"refreshTokenExpiredAt"`
}

func (a *authUseCase) Login(ctx context.Context, body *transport.LoginBody) (result *LoginResult, err error) {
	hashPassword := a.authSvc.HashPassword(body.Password, a.config.Secret.Password)
	user, err := a.userRepo.GetByEmailPassword(ctx, &ent.User{
		Email:    &body.Email,
		Password: &hashPassword,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, "이메일 혹은 비밀번호를 다시 확인해주세요.")
			return
		}
		return
	}

	if !user.IsEmailVerified {
		err = common.NewUnauthorizedError(ex.ErrUserEmailUnauthorization, nil)
		return
	}

	userType := ""

	if user.Type != nil {
		userType = *user.Type.ToString()
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, err := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.RefreshTokenExpired)
	if err != nil {
		return
	}

	if err = a.store.HSet(ctx, strconv.Itoa(user.ID), uid, strconv.Itoa(user.ID), time.Duration(a.config.JWT.RefreshTokenExpired)*time.Minute); err != nil {
		return
	}

	access, err := a.tokenSvc.GenerateToken(uid, user.ID, "normal", userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return
	}

	result = &LoginResult{
		AccessToken:           access,
		AccessTokenExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
		RefreshToken:          refresh,
		RefreshTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.RefreshTokenExpired),
	}

	return
}

func (a *authUseCase) SocialLogin(ctx context.Context, socialName model.SocialType, code string) (result *LoginResult, err error) {
	socialNameString := socialName.ToString()
	if socialNameString != &model.GoogleString && socialNameString != &model.KakaoString && socialNameString != &model.NaverString {
		err = errors.New("입력 값을 확인해주세요")
		return
	}

	user := new(ent.User)
	// 카카오 로그인
	if socialNameString == &model.KakaoString {
		kakaoInfo, errA := a.authSvc.GetKakaoInfo(code)
		if errA != nil {
			err = errA
			return
		}

		kakaoId := strconv.FormatUint(uint64(kakaoInfo.Id), 10)

		user.SocialKey = &kakaoId
		user.Email = kakaoInfo.KakaoAccount.Email
		user.Nickname = &kakaoInfo.KakaoAccount.Profile.NickName
		user.SocialName = &model.KakaoSocialType

		// 구글 로그인
	} else if socialNameString == &model.GoogleString {

		googleInfo, errA := a.authSvc.GetGoogleInfo(code)
		if errA != nil {
			err = errA
			return
		}

		googleId := googleInfo.Sub
		user.SocialKey = &googleId
		user.Email = &googleInfo.Email
		user.Nickname = &googleInfo.Nickname
		user.SocialName = &model.GoogleSocialType

		// 네이버 로그인
	} else if socialNameString == &model.NaverString {

		naverInfo, errA := a.authSvc.GetNaverInfo(code)
		if errA != nil {
			err = errA
			return
		}
		user.Nickname = &naverInfo.NickName
		user.SocialKey = &naverInfo.Id
		user.Email = &naverInfo.Email
		user.SocialName = &model.NaverSocialType

	}

	// 이미 존재하는 회원인지 확인.
	checkedUser, err := a.userRepo.GetBySocialKey(ctx, user)
	if err != nil && !ent.IsNotFound(err) {
		return
	}

	// 유저가 없으면 회원 정보 생성
	if checkedUser == nil {
		user.TermAgreeAt = time.Now()
		checkedUser, err = a.userRepo.Create(ctx, user)
		if err != nil {
			return
		}
	}

	userType := ""
	if checkedUser.Type != nil {
		userType = *user.Type.ToString()
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, err := a.tokenSvc.GenerateToken(uid, user.ID, *socialNameString, userType, a.config.JWT.RefreshTokenExpired)
	if err != nil {
		return
	}

	if err = a.store.HSet(ctx, strconv.Itoa(user.ID), uid, strconv.Itoa(user.ID), time.Duration(a.config.JWT.RefreshTokenExpired)*time.Minute); err != nil {
		return
	}

	access, err := a.tokenSvc.GenerateToken(uid, user.ID, *socialNameString, userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return
	}

	result = &LoginResult{
		AccessToken:           access,
		AccessTokenExpiredAt:  a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
		RefreshToken:          refresh,
		RefreshTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.RefreshTokenExpired),
	}
	return
}

func (a *authUseCase) GoogleRedirectUrl(ctx context.Context) string {
	return a.authSvc.GetGoogleRedirectUrl()
}

func (a *authUseCase) SocialSignUp(ctx context.Context, body *transport.SocialSignUpBody) (err error) {
	u, err := a.userRepo.Get(ctx, body.UserID)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return
	}

	if u.Email != nil {
		err = ex.NewConflictError(ex.ErrUserEmailAlreadyRegisted, nil)
		return
	}

	if err = a.userRepo.UpdateEmail(ctx, body.Email, body.UserID); err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrUserEmailAlreadyExist, nil)
		}
		return
	}
	return
}

func (a *authUseCase) SignUp(ctx context.Context, body *transport.SignUpBody) (err error) {
	hashPassword := a.authSvc.HashPassword(body.Password, a.config.Secret.Password)

	_, err = a.userRepo.Create(ctx, &ent.User{
		Email:       &body.Email,
		Password:    &hashPassword,
		Nickname:    &body.NickName,
		TermAgreeAt: time.Now(),
	})
	if err != nil {
		if ent.IsConstraintError(err) {
			err = ex.NewConflictError(ex.ErrUserEmailAlreadyExist, nil)
			return
		}
		return
	}

	// 이메일 인증을 위한 키 store에 저장.
	key := a.authSvc.GenerateRandomString()
	if err = a.store.Set(ctx, body.Email, key, time.Duration(time.Hour*24)); err != nil {
		return
	}

	go a.authSvc.SendEmailVerifiedUser(body.Email, key, time.Now().Format(time.RFC3339), a.config.Onthemat.HOST)

	return
}

func (a *authUseCase) CheckDuplicatedEmail(ctx context.Context, email string) (err error) {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	if isExist {
		err = ex.NewConflictError(ex.ErrUserEmailAlreadyExist, nil)
		return
	}
	return
}

func (a *authUseCase) VerifiyEmail(ctx context.Context, email string, authKey string, issuedAt string) (err error) {
	if a.authSvc.IsExpiredEmailForVerify(issuedAt) {
		err = ex.NewAuthenticationFailedError(ex.ErrEmailForVerifyExpired, nil)
		return
	}

	key := a.store.Get(ctx, email)

	if key != authKey {
		err = ex.NewBadRequestError(ex.ErrRandomKeyForEmailVerfiyUnavailable, nil)
		return
	}

	u, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return
	}

	if u.IsEmailVerified {
		err = ex.NewConflictError(ex.ErrUserEmailAlreadyVerfied, nil)
		return
	}

	if err = a.userRepo.UpdateEmailVerifeid(ctx, u.ID); err != nil {
		err = ex.NewConflictError(ex.ErrUserEmailAlreadyVerfied, nil)
		return
	}

	if err = a.store.Del(ctx, email); err != nil {
		err = ex.NewConflictError(ex.ErrUserEmailAlreadyVerfied, nil)
		return
	}

	return
}

func (a *authUseCase) SendEmailResetPassword(ctx context.Context, email string) (err error) {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return
	}

	if !isExist {
		err = ex.NewNotFoundError(ex.ErrUserNotFound, "존재하지 않는 이메일입니다.")
		return
	}

	hashPassword := a.authSvc.HashPassword(a.authSvc.GenerateRandomPassword(), a.config.Secret.Password)

	u := &ent.User{
		Email:        &email,
		TempPassword: &hashPassword,
	}
	if err = a.userRepo.UpdateTempPassword(ctx, u); err != nil {
		return
	}

	go a.authSvc.SendEmailResetPassword(u)

	return
}

type RefreshResult struct {
	AccessToken          string    `json:"accessToken"`
	AccessTokenExpiredAt time.Time `json:"accessTokenExpiredAt"`
}

func (a *authUseCase) Refresh(ctx context.Context, authorizationHeader []byte) (result *RefreshResult, err error) {
	refreshToken, err := a.authSvc.ExtractTokenFromHeader(string(authorizationHeader))
	if err != nil {
		err = ex.NewBadRequestError(ex.ErrAuthorizationHeaderFormatUnavailable, "Bearer")
		return
	}

	claim := new(token.TokenClaim)
	if err = a.tokenSvc.ParseToken(refreshToken, claim); err != nil {
		if err.Error() == jwt.ErrExiredToken {
			err = ex.NewUnauthorizedError(ex.ErrTokenExpired, nil)
			return
		}

		if err.Error() == jwt.ErrInvalidToken {
			err = ex.NewBadRequestError(ex.ErrTokenInvalid, nil)
			return
		}
		return
	}

	userIdString := strconv.Itoa(claim.UserId)
	val, err := a.store.HGet(ctx, userIdString, claim.Uuid)
	if err != nil {
		return
	}
	if val != userIdString {
		err = ex.NewBadRequestError(ex.ErrTokenInvalid, nil)
		return
	}

	u, err := a.userRepo.Get(ctx, claim.UserId)
	if err != nil {
		if ent.IsNotFound(err) {
			err = ex.NewNotFoundError(ex.ErrUserNotFound, nil)
			return
		}
		return

	}

	uid := uuid.New().String()

	userType := ""
	if u.Type != nil {
		userType = *u.Type.ToString()
	}

	loginType := "normal"
	if u.SocialName != nil {
		loginType = *u.SocialName.ToString()
	}

	access, err := a.tokenSvc.GenerateToken(uid, u.ID, loginType, userType, a.config.JWT.AccessTokenExpired)
	if err != nil {
		return
	}

	result = &RefreshResult{
		AccessToken:          access,
		AccessTokenExpiredAt: a.tokenSvc.GetExpiredAt(a.config.JWT.AccessTokenExpired),
	}
	return
}
