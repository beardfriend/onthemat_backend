package usecase_test

import (
	"context"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/ent"
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
