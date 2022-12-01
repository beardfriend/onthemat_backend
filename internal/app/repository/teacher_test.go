package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"
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
	areaRepo    AreaRepository
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
	ts.areaRepo = NewAreaRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *TeacherRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
}

// 각 테스트 종료 후 N회
func (ts *TeacherRepositoryTestSuite) TearDownTest() {
	// 모든 데이터 지우기
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

func (ts *TeacherRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if testName == "TestPatch" {
		ts.initData()
		u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
		ts.userNo = u.ID
		err := ts.teacherRepo.Create(ts.ctx, &ent.Teacher{
			UserID: u.ID,
			Name:   "Name",
			Age:    utils.Int(21),
			Edges: ent.TeacherEdges{
				Yoga: []*ent.Yoga{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				},
				WorkExperience: []*ent.TeacherWorkExperience{
					{
						AcademyName: "하타요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "아쉬탕가요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "하스",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "미무",
						WorkStartAt: transport.TimeString(time.Now()),
					},
				},
				Certification: []*ent.TeacherCertification{
					{
						AgencyName:   "agency",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
					{
						AgencyName:   "agency2",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
				},
				YogaRaw: []*ent.YogaRaw{
					{
						Name:      "아쉬탕가",
						TeacherID: utils.Int(1),
					},
					{
						Name:      "하타",
						TeacherID: utils.Int(1),
					},
				},
				Sigungu: []*ent.AreaSiGungu{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				},
			},
		})
		ts.NoError(err)
	}
	if testName == "TestUpdate" {
		ts.initData()
		u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
		ts.userNo = u.ID
		data := &ent.Teacher{
			UserID: u.ID,
			Name:   "Name",
			Age:    utils.Int(21),
			Edges: ent.TeacherEdges{
				Yoga: []*ent.Yoga{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				},
				WorkExperience: []*ent.TeacherWorkExperience{
					{
						AcademyName: "하타요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "아쉬탕가요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "하스",
						WorkStartAt: transport.TimeString(time.Now()),
					},
					{
						AcademyName: "미무",
						WorkStartAt: transport.TimeString(time.Now()),
					},
				},
				Certification: []*ent.TeacherCertification{
					{
						AgencyName:   "agency",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
					{
						AgencyName:   "agency2",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
				},
				YogaRaw: []*ent.YogaRaw{
					{
						Name:      "아쉬탕가",
						TeacherID: utils.Int(1),
					},
					{
						Name:      "하타",
						TeacherID: utils.Int(1),
					},
				},
				Sigungu: []*ent.AreaSiGungu{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				},
			},
		}
		err := ts.teacherRepo.Create(ts.ctx, data)
		// forPostMan, _ := json.Marshal(data)
		// fmt.Println(string(forPostMan))
		ts.NoError(err)
	}
}

func (ts *TeacherRepositoryTestSuite) TestUpdate() {
	ts.Run("성공", func() {
		err := ts.teacherRepo.Update(ts.ctx, &ent.Teacher{
			ID:     1,
			Name:   "name",
			UserID: 1,
			Edges: ent.TeacherEdges{
				WorkExperience: []*ent.TeacherWorkExperience{
					{
						ID:          1,
						AcademyName: "변경이름",
						WorkStartAt: transport.TimeString(time.Now()),
					},

					{
						ID:          3,
						AcademyName: "아쉬탕가요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},

					{
						ID:          10,
						AcademyName: "아쉬탕가요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},

					{
						ID:          11,
						AcademyName: "아쉬탕가요가원",
						WorkStartAt: transport.TimeString(time.Now()),
					},
				},
				Certification: []*ent.TeacherCertification{
					{
						ID:           1,
						AgencyName:   "change",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
					{
						ID:           4,
						AgencyName:   "added",
						TeacherID:    1,
						ClassStartAt: transport.TimeString(time.Now()),
					},
				},
				YogaRaw: []*ent.YogaRaw{
					{
						ID:        1,
						Name:      "chnaged",
						TeacherID: utils.Int(1),
					},
					{
						ID:        4,
						Name:      "added",
						TeacherID: utils.Int(1),
					},
				},
			},
		})
		ts.NoError(err)
	})
}

func (ts *TeacherRepositoryTestSuite) TestPatch() {
	ts.Run("성공", func() {
		ti := transport.TimeString(time.Now())
		err := ts.teacherRepo.Patch(ts.ctx, &request.TeacherPatchBody{
			WorkExperiences: []*request.WorkExperiencesForPatch{
				{
					Id:          utils.Int(1),
					AcademyName: utils.String("변경이름"),
					WorkStartAt: &ti,
				},

				{
					Id:          utils.Int(3),
					AcademyName: utils.String("수정"),
					WorkStartAt: &ti,
				},

				{
					AcademyName: utils.String("생성"),
					WorkStartAt: &ti,
				},

				{
					AcademyName: utils.String("생성"),
					WorkStartAt: &ti,
				},
			},
		}, 1, 1)
		ts.NoError(err)
	})
}

func TestTeacherRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TeacherRepositoryTestSuite))
}

func (ts *TeacherRepositoryTestSuite) initData() {
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

	bulk := make([]*ent.YogaGroupCreate, 4)
	bulk[0] = ts.client.YogaGroup.Create().SetCategory("아쉬탕가").SetCategoryEng("ashtanga").SetDescription("아쉬탕가 요가입니다.")
	bulk[1] = ts.client.YogaGroup.Create().SetCategory("하타").SetCategoryEng("hatha").SetDescription("하타 요가입니다.")
	bulk[2] = ts.client.YogaGroup.Create().SetCategory("빈야사").SetCategoryEng("vinyasa").SetDescription("빈야사 요가입니다.")
	bulk[3] = ts.client.YogaGroup.Create().SetCategory("아디다스").SetCategoryEng("adidas").SetDescription("아디다스 요가입니다.")
	err := ts.client.YogaGroup.CreateBulk(bulk...).Exec(context.Background())
	ts.NoError(err)
	bulks := make([]*ent.YogaCreate, 8)
	bulks[0] = ts.client.Yoga.Create().SetNameKor("아쉬탕가 레드").SetLevel(5).SetYogaGroupID(1)
	bulks[1] = ts.client.Yoga.Create().SetNameKor("아쉬탕가 프라이머리").SetLevel(4).SetYogaGroupID(1)
	bulks[2] = ts.client.Yoga.Create().SetNameKor("아쉬탕가 프라이머리 하프").SetLevel(3).SetYogaGroupID(1)
	bulks[3] = ts.client.Yoga.Create().SetNameKor("하타 플로우").SetLevel(3).SetYogaGroupID(2)
	bulks[4] = ts.client.Yoga.Create().SetNameKor("하타 테라피").SetLevel(2).SetYogaGroupID(2)
	bulks[5] = ts.client.Yoga.Create().SetNameKor("하타 기초").SetLevel(2).SetYogaGroupID(2)
	bulks[6] = ts.client.Yoga.Create().SetNameKor("빈야사 플로우").SetLevel(2).SetYogaGroupID(3)
	bulks[7] = ts.client.Yoga.Create().SetNameKor("빈야사 기초").SetLevel(1).SetYogaGroupID(3)
	err = ts.client.Yoga.CreateBulk(bulks...).Exec(context.Background())
	ts.NoError(err)
}
