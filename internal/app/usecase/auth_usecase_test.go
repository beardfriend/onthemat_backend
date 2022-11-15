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
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/ent"
	"onthemat/pkg/kakao"
	pkgMock "onthemat/pkg/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestAuthUC_TestSuite(t *testing.T) {
	suite.Run(t, new(AuthUC_TestSuite))
}

type AuthUC_TestSuite struct {
	suite.Suite
	authUC usecase.AuthUseCase

	mockTokenService *mocks.TokenService
	mockUserRepo     *mocks.UserRepository
	mockAuthService  *mocks.AuthService
	mockStore        *pkgMock.Store
}

// 모든 테스트 시작 전 1회
func (ts *AuthUC_TestSuite) SetupSuite() {
	c := config.NewConfig()

	ts.mockTokenService = new(mocks.TokenService)

	ts.mockUserRepo = new(mocks.UserRepository)
	ts.mockAuthService = new(mocks.AuthService)
	ts.mockStore = new(pkgMock.Store)
	ts.authUC = usecase.NewAuthUseCase(ts.mockTokenService, ts.mockUserRepo, ts.mockAuthService, ts.mockStore, c)
}

func (ts *AuthUC_TestSuite) TestAuthUC_CheckDuplicateEmail() {
	email := "sample@mail.com"
	ts.Run("success", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(false, nil).Once()
		err := ts.authUC.CheckDuplicatedEmail(context.TODO(), email)
		ts.NoError(err)
	})
	ts.Run("email-already-exisit", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(true, nil).Once()
		err := ts.authUC.CheckDuplicatedEmail(context.TODO(), email)
		errorStruct := err.(common.HttpError)
		ts.Equal(409, errorStruct.ErrCode)
		ts.Equal("이미 존재하는 이메일입니다.", errorStruct.ErrDetails)
	})
}

func (ts *AuthUC_TestSuite) TestAuthUC_Login() {
	userEmail := "asd@naver.com"
	userPassword := "password"

	ts.Run("유저 정보가 없을 경우", func() {
		ts.mockAuthService.On("HashPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("hashedPassword")
		ts.mockUserRepo.On("GetByEmailPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(nil, &ent.NotFoundError{}).Once()

		_, err := ts.authUC.Login(context.TODO(), &transport.LoginBody{
			Email:    "asd@naver.com",
			Password: "password",
		})

		errorStruct := err.(common.HttpError)
		ts.Equal(404, errorStruct.ErrCode)
		ts.Equal("이메일 혹은 비밀번호를 다시 확인해주세요.", errorStruct.ErrDetails)
	})

	ts.Run("이메일 인증이 되지 않은 경우", func() {
		ts.mockAuthService.On("HashPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("hashedPassword")
		ts.mockUserRepo.On("GetByEmailPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(&ent.User{
			Email:           &userEmail,
			Password:        &userPassword,
			IsEmailVerified: false,
		}, nil).Once()

		_, err := ts.authUC.Login(context.TODO(), &transport.LoginBody{
			Email:    userEmail,
			Password: userPassword,
		})

		errorStruct := err.(common.HttpError)
		ts.Equal(400, errorStruct.ErrCode)
		ts.Equal("이메일 인증이 필요합니다.", errorStruct.ErrDetails)
	})
}

func (ts *AuthUC_TestSuite) TestAuthUC_SocialLogin() {
	redirectCode := "examplecode"
	sociaKey := "123123123"
	email := "asd@naver.com"
	nickname := "nickname"

	ts.Run("성공", func() {
		ts.Run("카카오 로그인 (최초 접근)", func() {
			keyInt, _ := strconv.Atoi(sociaKey)
			ts.mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
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

			ts.mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(nil, nil).
				Once()

			ts.mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*ent.User")).
				Return(&ent.User{
					ID:          1,
					TermAgreeAt: time.Now(),
					SocialName:  &model.KakaoSocialType,
					SocialKey:   &sociaKey,
					Email:       &email,
					Nickname:    &nickname,
				}, nil).
				Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			ts.mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("AccessToken", nil).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).Once()

			// 검증
			l, err := ts.authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			ts.Equal(l.AccessToken, "AccessToken")
			ts.NoError(err, nil)
		})

		ts.Run("카카오 로그인 이미 가입한 유저", func() {
			keyInt, _ := strconv.Atoi(sociaKey)
			ts.mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
				Return(&kakao.GetUserInfoSuccessBody{
					Id: uint(keyInt),
				}, nil).
				Once()

			ts.mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(&ent.User{ID: 1, SocialKey: &sociaKey, SocialName: &model.KakaoSocialType}, nil).
				Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			ts.mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("AccessToken", nil).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).Once()

			l, err := ts.authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			ts.Equal(l.AccessToken, "AccessToken")
			ts.NoError(err, nil)
		})

		ts.Run("카카오 로그인 학원선생님 인증을 마친 유저", func() {
			keyInt, _ := strconv.Atoi(sociaKey)
			ts.mockAuthService.On("GetKakaoInfo", mock.AnythingOfType("string")).
				Return(&kakao.GetUserInfoSuccessBody{
					Id: uint(keyInt),
				}, nil).
				Once()

			ts.mockUserRepo.On("GetBySocialKey", mock.Anything, mock.Anything).
				Return(&ent.User{
					ID:         1,
					SocialKey:  &sociaKey,
					SocialName: &model.KakaoSocialType,
					Email:      &email,
					Nickname:   &nickname,
					Type:       &model.AcademyType,
				}, nil).
				Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("refreshToken", nil).
				Once()

			ts.mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()

			ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
				Return("AccessToken", nil).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).
				Once()

			ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
				Return(time.Now()).Once()

			l, err := ts.authUC.SocialLogin(context.TODO(), model.KakaoSocialType, redirectCode)

			ts.Equal(l.AccessToken, "AccessToken")
			ts.NoError(err, nil)
		})
	})
}
