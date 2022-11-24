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

type AcademyRepositoryTestSuite struct {
	suite.Suite
	config            *config.Config
	client            *ent.Client
	academyRepository AcademyRepository
	userRepo          UserRepository
	ctx               context.Context

	// data

	userId int
}

// 모든 테스트 시작 전 1회
func (ts *AcademyRepositoryTestSuite) SetupSuite() {
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
	ts.academyRepository = NewAcademyRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *AcademyRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *AcademyRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *AcademyRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "AcademyRepositoryTestSuite" {
		switch testName {
		}
	}
}

func (ts *AcademyRepositoryTestSuite) TestCreate() {
}

func TestAcademyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AcademyRepositoryTestSuite))
}
