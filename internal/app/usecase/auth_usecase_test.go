package usecase_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/model"
	"onthemat/internal/app/service/token"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/auth/jwt"
	"onthemat/pkg/ent"
	"onthemat/pkg/kakao"
	pkgMock "onthemat/pkg/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUC_CheckDuplicateEmail(t *testing.T) {
	c := config.NewConfig()

	mockTokenService := new(mocks.TokenService)

	mockUserRepo := new(mocks.UserRepository)
	mockAuthService := new(mocks.AuthService)
	mockStore := new(pkgMock.Store)

	email := "sample@mail.com"
	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(false, nil).Once()

		authUC := usecase.NewAuthUseCase(mockTokenService, mockUserRepo, mockAuthService, mockStore, c)
		err := authUC.CheckDuplicatedEmail(context.TODO(), email)
		assert.NoError(t, err)
	})

	t.Run("email-already-exisit", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(true, nil).Once()

		authUC := usecase.NewAuthUseCase(mockTokenService, mockUserRepo, mockAuthService, mockStore, c)
		err := authUC.CheckDuplicatedEmail(context.TODO(), email)
		errorStruct := err.(common.HttpError)
		assert.Equal(t, 409, errorStruct.ErrCode)
		assert.Equal(t, "이미 존재하는 이메일입니다.", errorStruct.ErrDetails)
	})
}

func TestAuthUC_Login(t *testing.T) {
	c := config.NewConfig()

	mockTokenService := new(mocks.TokenService)
	mockUserRepo := new(mocks.UserRepository)
	mockAuthService := new(mocks.AuthService)
	mockStore := new(pkgMock.Store)

	userEmail := "asd@naver.com"
	userPassword := "password"

	t.Run("유저 정보가 없을 경우", func(t *testing.T) {
		mockUserRepo.On("GetByEmailPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(nil, &ent.NotFoundError{}).Once()

		authUC := usecase.NewAuthUseCase(mockTokenService, mockUserRepo, mockAuthService, mockStore, c)
		_, err := authUC.Login(context.TODO(), &transport.LoginBody{
			Email:    "asd@naver.com",
			Password: "password",
		})
		errorStruct := err.(common.HttpError)
		assert.Equal(t, 404, errorStruct.ErrCode)
		assert.Equal(t, "이메일 혹은 비밀번호를 다시 확인해주세요.", errorStruct.ErrDetails)
	})

	t.Run("이메일 인증이 되지 않은 경우", func(t *testing.T) {
		mockUserRepo.On("GetByEmailPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(&ent.User{
			Email:           &userEmail,
			Password:        &userPassword,
			IsEmailVerified: false,
		}, nil).Once()

		authUC := usecase.NewAuthUseCase(mockTokenService, mockUserRepo, mockAuthService, mockStore, c)
		_, err := authUC.Login(context.TODO(), &transport.LoginBody{
			Email:    userEmail,
			Password: userPassword,
		})

		errorStruct := err.(common.HttpError)
		assert.Equal(t, 400, errorStruct.ErrCode)
		assert.Equal(t, "이메일 인증이 필요합니다.", errorStruct.ErrDetails)
	})
}

func TestAuthUC_SocialLogin(t *testing.T) {
	c := config.NewConfig()

	mockTokenService := new(mocks.TokenService)
	mockUserRepo := new(mocks.UserRepository)
	mockAuthService := new(mocks.AuthService)
	mockStore := new(pkgMock.Store)
	jwt := jwt.NewJwt().WithSignKey("signKey").Init()

	tokenSvc := token.NewToken(jwt)

	redirectCode := "examplecode"
	sociaKey := "123123123"
	email := "asd@naver.com"
	nickname := "nickname"

	t.Run("성공", func(t *testing.T) {
		t.Run("카카오 로그인 (최초 접근)", func(t *testing.T) {
			keyInt, _ := strconv.Atoi(sociaKey)
			mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
				Return(&kakao.GetUserInfoSuccessBody{
					Id: uint(keyInt),
					KakaoAccount: struct {
						Email   *string "json:\"email\""
						Profile struct {
							NickName string "json:\"nickname\""
						} "json:\"profile\""
					}{
						Email: &email,
						Profile: struct {
							NickName string "json:\"nickname\""
						}{
							NickName: nickname,
						},
					},
				}, nil).
				Once()

			mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(nil, nil).
				Once()

			mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*ent.User")).
				Return(&ent.User{
					ID:          1,
					TermAgreeAt: time.Now(),
					SocialName:  &model.KakaoSocialType,
					SocialKey:   &sociaKey,
					Email:       &email,
					Nickname:    &nickname,
				}, nil).
				Once()

			mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("AccessToken", nil).
				Once()

			mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).
				Once()

			mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).Once()

			// 검증
			authUC := usecase.NewAuthUseCase(mockTokenService, mockUserRepo, mockAuthService, mockStore, c)
			l, err := authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			assert.Equal(t, l.AccessToken, "AccessToken")
			assert.NoError(t, err, nil)
		})

		t.Run("카카오 로그인 이미 가입한 유저", func(t *testing.T) {
			keyInt, _ := strconv.Atoi(sociaKey)
			mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
				Return(&kakao.GetUserInfoSuccessBody{
					Id: uint(keyInt),
				}, nil).
				Once()

			mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(&ent.User{ID: 1, SocialKey: &sociaKey, SocialName: &model.KakaoSocialType}, nil).
				Once()

			mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("AccessToken", nil).
				Once()

			mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).
				Once()

			mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).Once()
			authUC := usecase.NewAuthUseCase(tokenSvc, mockUserRepo, mockAuthService, mockStore, c)
			l, err := authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			claim := &token.TokenClaim{}
			tokenSvc.ParseToken(l.AccessToken, claim)
			assert.Equal(t, claim.LoginType, "kakao")
			assert.Equal(t, claim.UserType, "")
			assert.NoError(t, err, nil)
		})

		t.Run("카카오 로그인 학원선생님 인증을 마친 유저", func(t *testing.T) {
			keyInt, _ := strconv.Atoi(sociaKey)
			mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
				Return(&kakao.GetUserInfoSuccessBody{
					Id: uint(keyInt),
				}, nil).
				Once()

			mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(&ent.User{
					ID:         1,
					SocialKey:  &sociaKey,
					SocialName: &model.KakaoSocialType,
					Email:      &email,
					Nickname:   &nickname,
					Type:       &model.AcademyType,
				}, nil).
				Once()

			mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			authUC := usecase.NewAuthUseCase(tokenSvc, mockUserRepo, mockAuthService, mockStore, c)
			l, err := authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			claim := &token.TokenClaim{}
			tokenSvc.ParseToken(l.AccessToken, claim)
			assert.Equal(t, claim.LoginType, "kakao")
			assert.Equal(t, claim.UserType, "academy")
			assert.NoError(t, err, nil)
		})
	})
}
