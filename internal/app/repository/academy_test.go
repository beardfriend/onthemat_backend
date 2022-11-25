package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	"github.com/goccy/go-json"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type AcademyRepositoryTestSuite struct {
	suite.Suite
	config            *config.Config
	client            *ent.Client
	academyRepository AcademyRepository
	userRepo          UserRepository
	areaRepo          AreaRepository
	ctx               context.Context

	// data
	create struct {
		userID int
	}
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
	ts.areaRepo = NewAreaRepository(ts.client)
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
		case "TestCreate":
			ts.initData()
		}
	}
}

func (ts *AcademyRepositoryTestSuite) TestCreate() {
	data := fmt.Sprintf(`
		{
			"user_id": %d,
			"sigungu_id": 1,
			"name": "name",
			"businessCode": "1138621886",
			"callNumber": "01064135418",
			"addressRoad": "서울시 양천구 도로명주소",
			"addressDetail": "상세주소"
		}
	`, ts.create.userID)

	yogaData := `
		[
			{
				"ID": 1
			},
			{
				"ID": 2
			}
		]
		
	`
	academy := &ent.Academy{}
	json.Unmarshal([]byte(data), academy)
	json.Unmarshal([]byte(yogaData), &academy.Edges.Yoga)
	err := ts.academyRepository.Create(ts.ctx, academy)
	ts.NoError(err)
}

func TestAcademyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AcademyRepositoryTestSuite))
}

func (ts *AcademyRepositoryTestSuite) initData() {
	u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
	ts.create.userID = u.ID

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
