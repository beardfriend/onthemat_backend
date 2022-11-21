package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	"github.com/stretchr/testify/suite"
)

type AreaRepositoryTestSuite struct {
	suite.Suite
	config   *config.Config
	client   *ent.Client
	areaRepo AreaRepository
	ctx      context.Context
}

// 모든 테스트 시작 전 1회
func (ts *AreaRepositoryTestSuite) SetupSuite() {
	t := ts.T()
	ts.ctx = context.Background()

	// 도커 디비 삭제 후 생성
	utils.RepoTestClose(t)
	time.Sleep(1 * time.Second)
	ts.config = utils.RepoTestInit(t)
	time.Sleep(3 * time.Second)
	// 포스트그레스 연결
	ts.client = infrastructure.NewPostgresDB(ts.config)

	// 모듈 연결
	ts.areaRepo = NewAreaRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *AreaRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *AreaRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *AreaRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "AreaRepositoryTestSuite" {
		if testName == "TestGetSigungu" {
			var sigungu []*ent.AreaSiGungu
			sigungu = append(sigungu, &ent.AreaSiGungu{
				Name:    "종로구",
				AdmCode: "11010",
			})
			sigungu = append(sigungu, &ent.AreaSiGungu{
				Name:    "중구",
				AdmCode: "11020",
			})

			e := ts.areaRepo.Create(ts.ctx, &ent.AreaSiDo{
				Name:    "서울특별자치시",
				AdmCode: "11",
				Version: 1,
			}, sigungu)
			ts.NoError(e)
		}
	}
}

func (ts *AreaRepositoryTestSuite) TestCreate() {
	ts.Run("성공", func() {
		var sigungu []*ent.AreaSiGungu
		sigungu = append(sigungu, &ent.AreaSiGungu{
			Name:    "종로구",
			AdmCode: "11010",
		})
		sigungu = append(sigungu, &ent.AreaSiGungu{
			Name:    "중구",
			AdmCode: "11020",
		})

		e := ts.areaRepo.Create(ts.ctx, &ent.AreaSiDo{
			Name:    "서울특별자치시",
			AdmCode: "11",
			Version: 1,
		}, sigungu)
		ts.NoError(e)

		sigungu, e = ts.client.AreaSiGungu.Query().All(ts.ctx)
		ts.NoError(e)
		ts.Equal(2, len(sigungu))
		si, e := ts.client.AreaSiDo.Query().All(ts.ctx)
		ts.NoError(e)
		ts.Equal(1, len(si))
	})
}

func (ts *AreaRepositoryTestSuite) TestGetSigungu() {
	data, err := ts.areaRepo.GetSigunGu(ts.ctx, "중구")
	ts.NoError(err)
	ts.Equal(data.Name, "중구")
}

func TestAreaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AreaRepositoryTestSuite))
}
