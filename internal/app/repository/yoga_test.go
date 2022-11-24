package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type YogaRepositoryTestSuite struct {
	suite.Suite
	config   *config.Config
	client   *ent.Client
	yogaRepo YogaRepository
	userRepo UserRepository
	ctx      context.Context

	testGetYogaRawData struct {
		userId int
	}

	testDeleteYogaRawData struct {
		userId int
		yogaId int
	}
}

// 모든 테스트 시작 전 1회
func (ts *YogaRepositoryTestSuite) SetupSuite() {
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
	ts.yogaRepo = NewYogaRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *YogaRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *YogaRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *YogaRepositoryTestSuite) TestCreate() {
	ts.Run("존재하지 않는 키 입력", func() {
		err := ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
			NameKor:     "아쉬탕가",
			NameEng:     nil,
			Description: nil,
			Level:       nil,
			Edges: ent.YogaEdges{
				YogaGroup: &ent.YogaGroup{
					ID: 1,
				},
			},
		})
		ts.Equal(ent.IsConstraintError(err), true)
	})
}

func TestYogaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(YogaRepositoryTestSuite))
}
