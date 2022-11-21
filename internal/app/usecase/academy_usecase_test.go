package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"onthemat/internal/app/common"
	"onthemat/internal/app/mocks"
	"onthemat/internal/app/model"
	"onthemat/internal/app/service"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/ent"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AcademyUCTestSuite struct {
	suite.Suite
	academyUC       usecase.AcademyUsecase
	mockAcademyRepo *mocks.AcademyRepository
	mockAcademySvc  *mocks.AcademyService
	mockUserRepo    *mocks.UserRepository
	mockAreaRepo    *mocks.AreaRepository
}

// 모든 테스트 시작 전 1회
func (ts *AcademyUCTestSuite) SetupSuite() {
	ts.mockAcademyRepo = new(mocks.AcademyRepository)
	ts.mockAcademySvc = new(mocks.AcademyService)
	ts.mockUserRepo = new(mocks.UserRepository)
	ts.mockAreaRepo = new(mocks.AreaRepository)

	ts.academyUC = usecase.NewAcademyUsecase(ts.mockAcademyRepo, ts.mockAcademySvc, ts.mockUserRepo, ts.mockAreaRepo)
}

// ------------------- Test Case -------------------

func (ts *AcademyUCTestSuite) TestCreate() {
	ts.Run("성공", func() {
		ts.mockAcademySvc.On("VerifyBusinessMan", mock.Anything).
			Return(nil).
			Once()

		ts.mockUserRepo.On("Get", mock.Anything, mock.Anything).
			Return(&ent.User{
				Type: nil,
			}, nil).Once()

		ts.mockAreaRepo.On("GetSigunGu", mock.Anything, mock.Anything).Return(&ent.AreaSiGungu{}, nil)

		ts.mockAcademyRepo.On("Create", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		err := ts.academyUC.Create(context.Background(), &transport.AcademyCreateRequestBody{}, 1)
		ts.NoError(err)
	})

	ts.Run("비즈니스 인증 실패 시", func() {
		ts.mockAcademySvc.On("VerifyBusinessMan", mock.Anything).
			Return(errors.New(service.ErrBussinessCodeInvalid)).
			Once()

		err := ts.academyUC.Create(context.Background(), &transport.AcademyCreateRequestBody{}, 1)
		errorStruct := err.(common.HttpError)
		ts.Equal(http.StatusBadRequest, errorStruct.ErrHttpCode)
	})

	ts.Run("이미 가입한 유저.", func() {
		ts.mockAcademySvc.On("VerifyBusinessMan", mock.Anything).
			Return(nil).
			Once()

		ts.mockUserRepo.On("Get", mock.Anything, mock.Anything).
			Return(&ent.User{
				Type: &model.TeacherType,
			}, nil).Once()

		err := ts.academyUC.Create(context.Background(), &transport.AcademyCreateRequestBody{}, 1)
		errorStruct := err.(common.HttpError)
		ts.Equal(http.StatusConflict, errorStruct.ErrHttpCode)
	})

	ts.Run("시군구 없을 때", func() {
	})
}

func (ts *AcademyUCTestSuite) TestGet() {
	ts.Run("성공", func() {
		ts.mockAreaRepo.On("GetSigunGu", mock.Anything, mock.Anything).Return(&ent.AreaSiGungu{}, nil)
		ts.mockAcademyRepo.On("Get", mock.Anything, mock.Anything).
			Return(&ent.Academy{}, nil).Once()

		_, err := ts.academyUC.Get(context.Background(), 1)
		ts.NoError(err)
	})

	ts.Run("NotFound", func() {
		ts.mockAcademyRepo.On("Get", mock.Anything, mock.Anything).
			Return(&ent.Academy{}, &ent.NotFoundError{}).Once()

		_, err := ts.academyUC.Get(context.Background(), 1)
		errorStruct := err.(common.HttpError)
		ts.Equal(http.StatusNotFound, errorStruct.ErrHttpCode)
	})
}

func (ts *AcademyUCTestSuite) TestUpdate() {
	ts.Run("성공", func() {
		ts.mockAreaRepo.On("GetSigunGu", mock.Anything, mock.Anything).Return(&ent.AreaSiGungu{}, nil)
		ts.mockAcademyRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()

		err := ts.academyUC.Update(context.Background(), &transport.AcademyUpdateRequestBody{}, 1)
		ts.NoError(err)
	})

	ts.Run("NotFound", func() {
		ts.mockAcademyRepo.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(&ent.NotFoundError{}).Once()

		err := ts.academyUC.Update(context.Background(), &transport.AcademyUpdateRequestBody{}, 1)
		errorStruct := err.(common.HttpError)
		ts.Equal(http.StatusNotFound, errorStruct.ErrHttpCode)
	})
	ts.Run("NotFound 시군구", func() {
	})
}

func (ts *AcademyUCTestSuite) TestList() {
	ts.mockAcademyRepo.On("Total", mock.Anything, mock.Anything, mock.Anything).
		Return(60, nil).Once()

	academies := make([]*ent.Academy, 10)
	ts.mockAcademyRepo.On("List", mock.Anything, mock.Anything, mock.Anything).
		Return(academies, nil).Once()

	_, p, e := ts.academyUC.List(context.Background(), &request.AcademyListQueries{
		PageNo:   1,
		PageSize: 10,
	})
	ts.Equal(p.PageCount, 6)
	ts.Equal(p.PageSize, 10)
	ts.Equal(p.RowCount, 10)
	fmt.Println(p)
	ts.NoError(e)
}

func TestAcademyUCTestSuite(t *testing.T) {
	suite.Run(t, new(AcademyUCTestSuite))
}
