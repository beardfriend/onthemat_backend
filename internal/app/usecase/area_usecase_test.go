package usecase_test

import (
	"context"
	"testing"

	"onthemat/internal/app/mocks"
	"onthemat/internal/app/service"
	"onthemat/internal/app/usecase"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AreaUsecaseTestSuite struct {
	suite.Suite

	mockAreaRepo    *mocks.AreaRepository
	mockAreaService *mocks.AreaService

	areaUsecase usecase.AreaUsecase
}

// 모든 테스트 시작 전 1회
func (ts *AreaUsecaseTestSuite) SetupSuite() {
	ts.mockAreaRepo = new(mocks.AreaRepository)
	ts.mockAreaService = new(mocks.AreaService)

	ts.areaUsecase = usecase.NewAreaUsecase(ts.mockAreaRepo, ts.mockAreaService)
}

// ------------------- Test Case -------------------

func (ts *AreaUsecaseTestSuite) TestCreate() {
	ts.mockAreaService.On("ParseExcelData", mock.Anything).
		Return([]service.Sido{{
			SidoName: "서울시",
			SidoCode: "11",
		}}, []service.Sigungu{{
			SigunguName: "양천구",
			SigunguCode: "110110",
		}, {
			SigunguName: "금천구",
			SigunguCode: "110110",
		}}, nil).
		Once()

	ts.mockAreaRepo.On("Create",
		mock.Anything,
		mock.AnythingOfType("*ent.AreaSiDo"),
		mock.AnythingOfType("[]*ent.AreaSiGungu"),
	).Return(nil).Once()

	err := ts.areaUsecase.CreateSiDo(context.Background(), "/fileurl")
	ts.NoError(err)
}

func TestAreaUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AreaUsecaseTestSuite))
}
