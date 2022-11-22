package usecase_test

import (
	"context"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/ent"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	userUsecase  usecase.UserUseCase
	mockUserRepo *mocks.UserRepository
}

// 모든 테스트 시작 전 1회
func (ts *UserUsecaseTestSuite) SetupSuite() {
	ts.mockUserRepo = new(mocks.UserRepository)
	ts.userUsecase = usecase.NewUserUseCase(ts.mockUserRepo)
}

// ------------------- Test Case -------------------

func (ts *UserUsecaseTestSuite) TestGetMe() {
	ts.Run("성공", func() {
		ts.mockUserRepo.On("Get", mock.Anything, mock.Anything).
			Return(&ent.User{
				ID: 1,
			}, nil).Once()

		user, err := ts.userUsecase.GetMe(context.Background(), 1)
		ts.NoError(err)
		ts.Equal(1, user.ID)
	})

	ts.Run("NotFound", func() {
		ts.mockUserRepo.On("Get", mock.Anything, mock.Anything).
			Return(nil, &ent.NotFoundError{}).Once()

		_, err := ts.userUsecase.GetMe(context.Background(), 1)
		errorStruct := err.(common.HttpError)
		ts.Equal(404, errorStruct.ErrHttpCode)
	})
}

func (ts *UserUsecaseTestSuite) TestAddYoga() {
	ts.Run("성공", func() {
		ts.mockUserRepo.On("AddYoga", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		ids := []int{1, 2}
		err := ts.userUsecase.AddYoga(context.Background(), 1, ids)
		ts.NoError(err)
	})

	ts.Run("Bad Request", func() {
		ts.mockUserRepo.On("AddYoga", mock.Anything, mock.Anything, mock.Anything).
			Return(&ent.ConstraintError{}).Once()

		ids := []int{1, 2}
		err := ts.userUsecase.AddYoga(context.Background(), 1, ids)
		errorStruct := err.(common.HttpError)
		ts.Equal(400, errorStruct.ErrHttpCode)
	})
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}
