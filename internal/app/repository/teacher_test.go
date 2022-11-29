package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/model"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type TeacherRepositoryTestSuite struct {
	suite.Suite
	config      *config.Config
	client      *ent.Client
	teacherRepo TeacherRepository
	userRepo    UserRepository
	ctx         context.Context

	// data

	userNo int
}

// 모든 테스트 시작 전 1회
func (ts *TeacherRepositoryTestSuite) SetupSuite() {
	t := ts.T()
	ts.ctx = context.Background()

	c := utils.GetTestConfig(t)
	ts.config = c
	// 포스트그레스 연결
	ts.client = infrastructure.NewPostgresDB(ts.config)

	// 모듈 연결
	ts.teacherRepo = NewTeacherRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *TeacherRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// 각 테스트 종료 후 N회
func (ts *TeacherRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *TeacherRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "TeacherRepositoryTestSuite" {
		if testName == "TestCreate" {
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.userNo = u.ID
		}
	}
}

func (ts *TeacherRepositoryTestSuite) TestCreate() {
	ts.Run("성공", func() {
		ts.teacherRepo.Create(ts.ctx, &ent.Teacher{
			Edges: ent.TeacherEdges{
				WorkExperience: []*ent.TeacherWorkExperience{
					{
						ClassContent: &[]model.ClassContent{
							{
								RunningTime: 1,
								YogaId:      1,
							},
						},
					},
				},
			},
		})
	})
}

func TestTeacherRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ImageRepositoryTestSuite))
}
