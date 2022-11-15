package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/model"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

func TestUserYogaRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserYogaRepositoryTestSuite))
}

type UserYogaRepositoryTestSuite struct {
	suite.Suite
	config       *config.Config
	client       *ent.Client
	userYogaRepo UserYogaRepository
	userRepo     UserRepository
	ctx          context.Context
	// Data
	exisitUserNo int
}

// 모든 테스트 시작 전 1회
func (ts *UserYogaRepositoryTestSuite) SetupSuite() {
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
	ts.userYogaRepo = NewUserYogaRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *UserYogaRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *UserYogaRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *UserYogaRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "UserYogaRepositoryTestSuite" {
		if testName == "TestCreate" {

			teacherType := model.TeacherType
			u, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Type: &teacherType,
			})
			ts.NoError(err)
			ts.exisitUserNo = u.ID

		}
	}
}

func (ts *UserYogaRepositoryTestSuite) TestCreate() {
	ts.Run("성공적인 생성", func() {
		var yoga []*ent.UserYoga
		// 1개
		teacherType := model.TeacherType
		yoga = append(yoga, &ent.UserYoga{
			Name:     "하타",
			UserType: teacherType,
		})
		// 2개
		yoga = append(yoga, &ent.UserYoga{
			Name:     "인",
			UserType: teacherType,
		})
		us, err := ts.userYogaRepo.CreateMany(ts.ctx, yoga, ts.exisitUserNo)
		ts.Equal(len(us), 2)
		ts.NoError(err)
	})

	ts.Run("존재하지 않는 유저 아이디로 생성하려는 경우", func() {
		var yoga []*ent.UserYoga
		// 1개
		yoga = append(yoga, &ent.UserYoga{
			Name:     "하타",
			UserType: model.TeacherType,
		})

		_, err := ts.userYogaRepo.CreateMany(ts.ctx, yoga, 100)
		ts.Equal(ent.IsConstraintError(err), true)
	})

	ts.Run("빈 데이터만 입력했을 때", func() {
		// 빈 배열을 넣어도 생성이 안되고 에러가 나지 않음.
		var yoga []*ent.UserYoga
		u, err := ts.userYogaRepo.CreateMany(ts.ctx, yoga, ts.exisitUserNo)
		ts.Equal(len(u), 0)
		ts.NoError(err)
	})
}
