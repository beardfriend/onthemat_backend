package usecase_test

import (
	"context"
	"fmt"
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
	"onthemat/pkg/ent"
	"onthemat/pkg/kakao"
	pkgMock "onthemat/pkg/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthUCTestSuite struct {
	suite.Suite
	authUC usecase.AuthUseCase

	mockTokenService *mocks.TokenService
	mockUserRepo     *mocks.UserRepository
	mockAuthService  *mocks.AuthService
	mockStore        *pkgMock.Store
	count            int
}

// 모든 테스트 시작 전 1회
func (ts *AuthUCTestSuite) SetupSuite() {
	c := config.NewConfig()

	ts.mockTokenService = new(mocks.TokenService)

	ts.mockUserRepo = new(mocks.UserRepository)
	ts.mockAuthService = new(mocks.AuthService)
	ts.mockStore = new(pkgMock.Store)
	ts.authUC = usecase.NewAuthUseCase(ts.mockTokenService, ts.mockUserRepo, ts.mockAuthService, ts.mockStore, c)
}

func (ts *AuthUCTestSuite) TearDownTest() {
	ts.count++
}

func (ts *AuthUCTestSuite) TearDownSuite() {
	fmt.Printf("-------------- %d 개 모듈-------------\n", ts.count)
}

// ------------------- Test Case -------------------

func (ts *AuthUCTestSuite) TestCheckDuplicateEmail() {
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
		ts.Equal(409, errorStruct.ErrHttpCode)
	})
}

func (ts *AuthUCTestSuite) TestSignUp() {
	ts.Run("이미 존재하는 이메일", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(true, nil).Once()
		err := ts.authUC.SignUp(context.TODO(), &transport.SignUpBody{
			Email:     "alreadyExisit@naver.com",
			Password:  "password",
			NickName:  "nick",
			TermAgree: true,
		})
		errorStruct := err.(common.HttpError)
		ts.Equal(409, errorStruct.ErrHttpCode)
	})

	ts.Run("회원가입 성공", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(false, nil).Once()
		ts.mockAuthService.On("HashPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("hashedPassword")
		ts.mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*ent.User")).Return(nil, nil).Once()
		ts.mockAuthService.On("GenerateRandomString").Return("randomasdqwd")
		ts.mockStore.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil).Once()
		ts.mockAuthService.On("SendEmailVerifiedUser", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
		err := ts.authUC.SignUp(context.TODO(), &transport.SignUpBody{
			Email:     "email@naver.com",
			Password:  "password",
			NickName:  "nick",
			TermAgree: true,
		})
		ts.NoError(err)
	})
}

func (ts *AuthUCTestSuite) TestVerifiyEmail() {
	ts.Run("인증 만료일이 지났을 때", func() {
		ts.mockAuthService.On("IsEmailForVerifyExpired", mock.AnythingOfType("string")).
			Return(true).Once()

		err := ts.authUC.VerifiyEmail(context.TODO(), "email@naver.com", "HackingRandomKey", "issuedAt")
		errorStruct := err.(common.HttpError)
		ts.Equal(401, errorStruct.ErrHttpCode)
	})
	ts.Run("인증키가 잘못됐을 때", func() {
		ts.mockAuthService.On("IsEmailForVerifyExpired", mock.AnythingOfType("string")).
			Return(false).Once()
		ts.mockStore.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return("randomasdqwd").
			Once()
		err := ts.authUC.VerifiyEmail(context.TODO(), "email@naver.com", "HackingRandomKey", "issuedAt")
		errorStruct := err.(common.HttpError)
		ts.Equal(400, errorStruct.ErrHttpCode)
	})

	ts.Run("이미 인증된 유저", func() {
		email := "email@naver.com"

		ts.mockAuthService.On("IsEmailForVerifyExpired", mock.AnythingOfType("string")).
			Return(false).Once()
		ts.mockStore.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return("randomasdqwd").
			Once()

		ts.mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(&ent.User{
				ID:              1,
				Email:           &email,
				IsEmailVerified: true,
			}, nil).Once()

		err := ts.authUC.VerifiyEmail(context.TODO(), email, "randomasdqwd", "issuedAt")
		errorStruct := err.(common.HttpError)
		ts.Equal(409, errorStruct.ErrHttpCode)
	})

	ts.Run("인증 성공", func() {
		email2 := "email2@naver.com"
		ts.mockAuthService.On("IsEmailForVerifyExpired", mock.AnythingOfType("string")).
			Return(false).Once()

		ts.mockStore.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return("randomasdqwd").
			Once()

		ts.mockUserRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(&ent.User{
				ID:              2,
				Email:           &email2,
				IsEmailVerified: false,
			}, nil).Once()

		ts.mockUserRepo.On("UpdateEmailVerifeid", mock.Anything, mock.AnythingOfType("int")).
			Return(nil).Once()

		ts.mockStore.On("Del", mock.Anything, mock.AnythingOfType("string")).
			Return(nil).
			Once()

		err := ts.authUC.VerifiyEmail(context.TODO(), email2, "randomasdqwd", "issuedAt")
		ts.NoError(err)
	})
}

func (ts *AuthUCTestSuite) TestSendEmailResetPassword() {
	userEmail := "asd@naver.com"

	ts.Run("존재하지 않는 이메일", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(false, nil).Once()

		err := ts.authUC.SendEmailResetPassword(context.TODO(), userEmail)

		errorStruct := err.(common.HttpError)
		ts.Equal(404, errorStruct.ErrHttpCode)
	})

	ts.Run("전송 성공", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(true, nil).Once()
		ts.mockAuthService.On("GenerateRandomPassword").Return("randomPassword").Once()
		ts.mockAuthService.On("HashPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("hashedPassword")
		ts.mockUserRepo.On("UpdateTempPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(nil).Once()
		ts.mockAuthService.On("SendEmailResetPassword", mock.AnythingOfType("*ent.User")).Return(nil).Once()
		err := ts.authUC.SendEmailResetPassword(context.TODO(), userEmail)

		ts.NoError(err)
	})
}

func (ts *AuthUCTestSuite) TestRefresh() {
	userEmail := "asd@naver.com"
	socialKey := "asadasd"
	ts.Run("성공 유저타입 [x] 소셜로그인[o]", func() {
		ts.mockAuthService.On("ExtractTokenFromHeader", mock.AnythingOfType("string")).
			Return("refreshToken", nil).Once()

		var claim token.TokenClaim
		ts.mockTokenService.On("ParseToken",
			mock.AnythingOfType("string"),
			&claim,
		).
			Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*token.TokenClaim)
			arg.UserId = 1
		}).Once()

		ts.mockStore.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return("1").Once()

		ts.mockUserRepo.On("Get", mock.Anything, mock.AnythingOfType("int")).
			Return(&ent.User{
				ID:         1,
				Email:      &userEmail,
				SocialKey:  &socialKey,
				SocialName: &model.KakaoSocialType,
			}, nil).Once()

		ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return("AccessToken", nil).
			Once()

		ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
			Return(time.Now()).
			Once()

		_, err := ts.authUC.Refresh(context.TODO(), []byte("Bearer refreshToken"))
		ts.NoError(err)
	})

	ts.Run("성공 유저타입 [o] 소셜로그인[o]", func() {
		ts.mockAuthService.On("ExtractTokenFromHeader", mock.AnythingOfType("string")).
			Return("refreshToken", nil).Once()

		var claim token.TokenClaim
		ts.mockTokenService.On("ParseToken",
			mock.AnythingOfType("string"),
			&claim,
		).
			Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*token.TokenClaim)
			arg.UserId = 1
		}).Once()

		ts.mockStore.On("Get", mock.Anything, mock.AnythingOfType("string")).
			Return("1").Once()

		ts.mockUserRepo.On("Get", mock.Anything, mock.AnythingOfType("int")).
			Return(&ent.User{
				ID:         1,
				Email:      &userEmail,
				SocialKey:  &socialKey,
				SocialName: &model.KakaoSocialType,
				Type:       &model.AcademyType,
			}, nil).Once()

		ts.mockTokenService.On("GenerateToken", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("int")).
			Return("AccessToken", nil).
			Once()

		ts.mockTokenService.On("GetExpiredAt", mock.AnythingOfType("int")).
			Return(time.Now()).
			Once()

		_, err := ts.authUC.Refresh(context.TODO(), []byte("Bearer refreshToken"))
		ts.NoError(err)
	})
}

func (ts *AuthUCTestSuite) TestLogin() {
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
		ts.Equal(404, errorStruct.ErrHttpCode)
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
		ts.Equal(401, errorStruct.ErrHttpCode)
	})

	ts.Run("성공", func() {
		ts.mockAuthService.On("HashPassword", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("hashedPassword")
		ts.mockUserRepo.On("GetByEmailPassword", mock.Anything, mock.AnythingOfType("*ent.User")).Return(&ent.User{
			Email:           &userEmail,
			Password:        &userPassword,
			Type:            &model.AcademyType,
			IsEmailVerified: true,
		}, nil).Once()

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
			Return(time.Now()).
			Once()

		l, err := ts.authUC.Login(context.TODO(), &transport.LoginBody{
			Email:    userEmail,
			Password: userPassword,
		})
		ts.NoError(err)
		ts.Equal(l.AccessToken, "AccessToken")
	})
}

func (ts *AuthUCTestSuite) TestSocialLogin() {
	redirectCode := "examplecode"
	sociaKey := "123123123"
	email := "asd@naver.com"
	nickname := "nickname"

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
}

func (ts *AuthUCTestSuite) TestSocialSignup() {
	ts.Run("email-already Exisit", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(true, nil).Once()
		err := ts.authUC.SocialSignUp(context.TODO(), &transport.SocialSignUpBody{
			UserID: 1,
			Email:  "asd@naver.com",
		})
		errorStruct := err.(common.HttpError)
		ts.Equal(409, errorStruct.ErrHttpCode)
	})
	ts.Run("success", func() {
		ts.mockUserRepo.On("FindByEmail", mock.Anything, mock.AnythingOfType("string")).Return(false, nil).Once()
		ts.mockUserRepo.On("UpdateEmail",
			mock.Anything,
			mock.AnythingOfType("string"),
			mock.AnythingOfType("int")).
			Return(nil).Once()

		err := ts.authUC.SocialSignUp(context.TODO(), &transport.SocialSignUpBody{
			UserID: 1,
			Email:  "asd@naver.com",
		})

		ts.NoError(err)
	})
}

func TestAuthUCTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUCTestSuite))
}
