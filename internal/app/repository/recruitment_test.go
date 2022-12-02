package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	"github.com/stretchr/testify/suite"
)

type RecruitmentTestSuite struct {
	suite.Suite
	config      *config.Config
	client      *ent.Client
	teacherRepo TeacherRepository
	userRepo    UserRepository
	areaRepo    AreaRepository
	recruitRepo RecruitmentRepository
	ctx         context.Context

	// data

	userNo int
}

func (ts *RecruitmentTestSuite) SetupSuite() {
	t := ts.T()
	ts.ctx = context.Background()

	c := utils.GetTestConfig(t)
	ts.config = c
	// 포스트그레스 연결
	ts.client = infrastructure.NewPostgresDB(ts.config)

	// 모듈 연결
	ts.teacherRepo = NewTeacherRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
	ts.areaRepo = NewAreaRepository(ts.client)
	ts.recruitRepo = NewRecruitmentRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *RecruitmentTestSuite) TearDownSuite() {
	ts.client.Close()
}

// 각 테스트 종료 후 N회
func (ts *RecruitmentTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *RecruitmentTestSuite) BeforeTest(suiteName, testName string) {
}

// raw query Test
func TestRecruitmentList(t *testing.T) {
	c := config.NewConfig()
	err := c.Load("../../../configs")
	if err != nil {
		t.Error(err)
	}

	db := infrastructure.NewPostgresDB(c)

	repo := NewRecruitmentRepository(db)
	ctx := context.Background()
	module := utils.NewPagination(1, 10)
	startTime, _ := time.Parse("2006-01-02T15:04:05", "2022-12-01T00:00:00")
	endTime, _ := time.Parse("2006-01-02T15:04:05", "2022-12-03T00:00:00")
	f := transport.TimeString(startTime)
	e := transport.TimeString(endTime)
	l, _ := repo.List(ctx, module, &f, &e, nil, nil)
	fmt.Println(l)
}
