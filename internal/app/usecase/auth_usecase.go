package usecase

import (
	"context"
	"errors"
	"strconv"
	"time"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/model"
	"onthemat/internal/app/repository"
	"onthemat/internal/app/service"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
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

	SendEmailResetPassword(ctx context.Context, email string) error
	CheckDuplicatedEmail(ctx context.Context, email string) error
	VerifiedEmail(ctx context.Context, email string, authKey string) error
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

func (a *authUseCase) Login(ctx context.Context, body *transport.LoginBody) (*LoginResult, error) {
	defer ctx.Done()
	hashPassword := a.authSvc.HashPassword(body.Password, a.config.Secret.Password)

	user, err := a.userRepo.GetByEmailPassword(ctx, &ent.User{
		Email:    &body.Email,
		Password: &hashPassword,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, common.NewNotFoundError("이메일 혹은 비밀번호를 다시 확인해주세요.")
		}
		return nil, err
	}

	if !user.IsEmailVerified {
		return nil, common.NewBadRequestError("이메일 인증이 필요합니다.")
	}

	userType := ""
	if user.Type == &model.AcademyType || user.Type == &model.TeacherType {
		userType = *user.Type.ToString()
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

func (a *authUseCase) SocialLogin(ctx context.Context, socialName model.SocialType, code string) (*LoginResult, error) {
	user := new(ent.User)
	socialNameString := socialName.ToString()

	if socialNameString != &model.GoogleString && socialNameString != &model.KakaoString && socialNameString != &model.NaverString {
		return nil, errors.New("입력 값을 확인해주세요")
	}

	// 카카오 로그인
	if socialNameString == &model.KakaoString {

		kakaoInfo, err := a.authSvc.GetKakaoInfo(code)
		if err != nil {
			return nil, err
		}

		kakaoId := strconv.FormatUint(uint64(kakaoInfo.Id), 10)

		user.SocialKey = &kakaoId
		user.Email = kakaoInfo.KakaoAccount.Email
		user.Nickname = &kakaoInfo.KakaoAccount.Profile.NickName
		user.SocialName = &model.KakaoSocialType

		// 구글 로그인
	} else if socialNameString == &model.GoogleString {

		googleInfo, err := a.authSvc.GetGoogleInfo(code)
		if err != nil {
			return nil, err
		}

		googleId := googleInfo.Sub
		user.SocialKey = &googleId
		user.Email = &googleInfo.Email
		user.Nickname = &googleInfo.Nickname
		user.SocialName = &model.GoogleSocialType

		// 네이버 로그인
	} else if socialNameString == &model.NaverString {

		naverInfo, err := a.authSvc.GetNaverInfo(code)
		if err != nil {
			return nil, err
		}
		user.Nickname = &naverInfo.NickName
		user.SocialKey = &naverInfo.Id
		user.Email = &naverInfo.Email
		user.SocialName = &model.NaverSocialType

	}

	// 이미 존재하는 회원인지 확인.
	checkedUser, err := a.userRepo.GetBySocialKey(ctx, user)
	if err != nil && !ent.IsNotFound(err) {
		return nil, err
	}

	// 유저가 없으면 회원 정보 생성
	if checkedUser == nil {
		user.TermAgreeAt = time.Now()
		checkedUser, err = a.userRepo.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	userType := ""
	if checkedUser.Type == &model.AcademyType || checkedUser.Type == &model.TeacherType {
		userType = *checkedUser.Type.ToString()
	}

	// 토큰 발행
	uid := uuid.New().String()
	refresh, err := a.tokenSvc.GenerateToken(uid, user.ID, *socialNameString, userType, a.config.JWT.RefreshTokenExpired)
	if err != nil {
		return nil, err
	}

	if err := a.store.Set(ctx, uid, strconv.Itoa(user.ID), time.Duration(a.config.JWT.RefreshTokenExpired)*time.Minute); err != nil {
		return nil, err
	}

	access, err := a.tokenSvc.GenerateToken(uid, user.ID, *socialNameString, userType, a.config.JWT.AccessTokenExpired)
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

func (a *authUseCase) GoogleRedirectUrl(ctx context.Context) string {
	return a.authSvc.GetGoogleRedirectUrl()
}

func (a *authUseCase) SocialSignUp(ctx context.Context, body *transport.SocialSignUpBody) error {
	_, err := a.userRepo.Update(ctx, &ent.User{
		ID:          body.UserID,
		Email:       &body.Email,
		TermAgreeAt: time.Now(),
		Type:        nil,
	})
	return err
}

func (a *authUseCase) SignUp(ctx context.Context, body *transport.SignUpBody) error {
	isExist, err := a.userRepo.FindByEmail(ctx, body.Email)
	if err != nil {
		return err
	}

	if isExist {
		return common.NewConflictError("이미 존재하는 이메일입니다.")
	}

	hashPassword := a.authSvc.HashPassword(body.Password, a.config.Secret.Password)

	_, err = a.userRepo.Create(ctx, &ent.User{
		Email:       &body.Email,
		Password:    &hashPassword,
		Nickname:    &body.NickName,
		TermAgreeAt: time.Now(),
	})

	if err != nil {
		return err
	}

	// 이메일 인증을 위한 키 store에 저장.
	key := a.authSvc.GenerateRandomString()
	if err := a.store.Set(ctx, body.Email, key, time.Duration(time.Hour*24)); err != nil {
		return err
	}

	go a.authSvc.SendEmailVerifiedUser(body.Email, key, a.config.Onthemat.HOST)

	return err
}

func (a *authUseCase) CheckDuplicatedEmail(ctx context.Context, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if isExist {
		return common.NewConflictError("이미 존재하는 이메일입니다.")
	}

	return nil
}

func (a *authUseCase) VerifiedEmail(ctx context.Context, email string, authKey string) error {
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

func (a *authUseCase) SendEmailResetPassword(ctx context.Context, email string) error {
	isExist, err := a.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if !isExist {
		return common.NewBadRequestError("존재하지 않는 이메일입니다.")
	}

	hashPassword := a.authSvc.HashPassword(a.authSvc.GenerateRandomPassword(), a.config.Secret.Password)

	u := &ent.User{
		Email:        &email,
		TempPassword: &hashPassword,
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

func (a *authUseCase) Refresh(ctx context.Context, authorizationHeader []byte) (*RefreshResult, error) {
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
		userType = *u.Type.ToString()
	}

	loginType := "normal"
	if u.SocialName != nil {
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
