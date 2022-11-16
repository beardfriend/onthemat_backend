package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/image"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type ImageRepositoryTestSuite struct {
	suite.Suite
	config          *config.Config
	client          *ent.Client
	imageRepository ImageRepository
	userRepo        UserRepository
	ctx             context.Context

	// data

	userNo int
}

// 모든 테스트 시작 전 1회
func (ts *ImageRepositoryTestSuite) SetupSuite() {
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
	ts.imageRepository = NewImageRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *ImageRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *ImageRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *ImageRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "ImageRepositoryTestSuite" {
		if testName == "TestCreate" {
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.userNo = u.ID
		}
	}
}

func (ts *ImageRepositoryTestSuite) TestCreate() {
	ts.Run("성공", func() {
		e := ts.imageRepository.Create(ts.ctx, &ent.Image{
			Name:        "name",
			Path:        "http://www.anver.com",
			Size:        1234123,
			ContentType: "asf/sdsd",
			Type:        image.Type("profile"),
		}, ts.userNo)
		ts.NoError(e)
	})
}

func TestImageRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ImageRepositoryTestSuite))
}
