package usecase_test

import (
	"context"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/usecase"
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
