package repository

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"

	fake "github.com/brianvoe/gofakeit/v6"
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

	c := utils.GetTestConfig(t)
	ts.config = c
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
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.create.userID = u.ID
			ts.initData()
			ts.userRepo.Create(ts.ctx, &ent.User{})
		case "TestUpdate":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.create.userID = u.ID
			ts.initData()
			err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
				UserID:        1,
				SigunguID:     1,
				Name:          "하타요가원",
				BusinessCode:  "1234",
				CallNumber:    "010122222222",
				AddressRoad:   "전체주소",
				AddressDetail: utils.String("상세주소"),
				Edges: ent.AcademyEdges{
					Yoga: []*ent.Yoga{
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
		case "TestGet":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.create.userID = u.ID
			ts.initData()
			err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
				UserID:        1,
				SigunguID:     1,
				Name:          "하타요가원",
				BusinessCode:  "1234",
				CallNumber:    "010122222222",
				AddressRoad:   "전체주소",
				AddressDetail: utils.String("상세주소"),
				Edges: ent.AcademyEdges{
					Yoga: []*ent.Yoga{
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
		case "TestPatch":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.create.userID = u.ID
			ts.initData()
			err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
				UserID:        1,
				SigunguID:     1,
				Name:          "하타요가원",
				BusinessCode:  "1234",
				CallNumber:    "010122222222",
				AddressRoad:   "전체주소",
				AddressDetail: utils.String("상세주소"),
				Edges: ent.AcademyEdges{
					Yoga: []*ent.Yoga{
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

		case "TestList":
			ts.initData()
			ts.initForList()

		}
	}
}

func (ts *AcademyRepositoryTestSuite) TestCreate() {
	ts.Run("with yoga", func() {
		data := fmt.Sprintf(`
		{
			"userId": %d,
			"sigunguId": 1,
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
	})

	ts.Run("without yoga", func() {
		err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
			UserID:        2,
			SigunguID:     1,
			Name:          "하타요가원",
			BusinessCode:  "1234",
			CallNumber:    "010122222222",
			AddressRoad:   "전체주소",
			AddressDetail: utils.String("상세주소"),
		})
		ts.NoError(err)
	})
}

func (ts *AcademyRepositoryTestSuite) TestUpdate() {
	ts.Run("success", func() {
		err := ts.academyRepository.Update(ts.ctx, &ent.Academy{
			ID:            1,
			UserID:        1,
			SigunguID:     2,
			Name:          "바꾸고",
			BusinessCode:  "1234",
			CallNumber:    "010122222222",
			AddressRoad:   "전체주소",
			AddressDetail: nil,
		})
		ts.NoError(err)

		academy, _ := ts.client.Academy.Get(ts.ctx, 1)
		ts.Equal("바꾸고", academy.Name)
		ts.Equal(0, len(academy.Edges.Yoga))
		ts.Nil(academy.AddressDetail)
	})

	ts.Run("success withYoga", func() {
		err := ts.academyRepository.Update(ts.ctx, &ent.Academy{
			ID:            1,
			UserID:        1,
			SigunguID:     2,
			Name:          "바꾸고",
			BusinessCode:  "1234",
			CallNumber:    "010122222222",
			AddressRoad:   "전체주소",
			AddressDetail: nil,
			Edges: ent.AcademyEdges{
				Yoga: []*ent.Yoga{
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

		academy, _ := ts.client.Academy.Get(ts.ctx, 1)
		ts.Equal("바꾸고", academy.Name)
		ts.Nil(academy.AddressDetail)
	})
}

func (ts *AcademyRepositoryTestSuite) TestPatch() {
	err := ts.academyRepository.Patch(ts.ctx, &request.AcademyPatchBody{
		Info: &request.AcademyInfoForPatch{
			Name: utils.String("학원이름"),
		},
	}, 1, 1)
	academy, _ := ts.client.Academy.Get(ts.ctx, 1)

	// Validate
	ts.Equal("학원이름", academy.Name)
	ts.Equal("010122222222", academy.CallNumber)
	ts.NoError(err)
}

func (ts *AcademyRepositoryTestSuite) TestGet() {
	ts.Run("success", func() {
		academy, err := ts.academyRepository.Get(ts.ctx, 1)

		// Validate
		ts.NoError(err)
		fmt.Println(academy)
	})
}

func (ts *AcademyRepositoryTestSuite) TestList() {
	ts.Run("DESC", func() {
		pageModule := utils.NewPagination(1, 10)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, nil, common.DESC)

		// Validate
		ts.Equal(10, len(list))
		ts.Greater(list[0].ID, list[1].ID)
		ts.NoError(err)
	})

	ts.Run("ASC", func() {
		pageModule := utils.NewPagination(1, 10)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, nil, common.ASC)

		// Validate
		ts.Equal(10, len(list))
		ts.Greater(list[1].ID, list[0].ID)
		ts.NoError(err)
	})

	ts.Run("ASC BY Name", func() {
		pageModule := utils.NewPagination(1, 5)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, utils.String("NAME"), common.ASC)
		nameFirstWordOfListZero := list[0].Name[0]
		nameFirstWordOfListOne := list[1].Name[0]
		firstInt, _ := strconv.Atoi(string(nameFirstWordOfListZero))
		secondInt, _ := strconv.Atoi(string(nameFirstWordOfListOne))

		// Validate
		ts.GreaterOrEqual(firstInt, secondInt)
		ts.NoError(err)
	})

	ts.Run("DESC BY Name", func() {
		pageModule := utils.NewPagination(1, 5)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, utils.String("NAME"), common.DESC)
		nameFirstWordOfListZero := list[0].Name[0]
		nameFirstWordOfListOne := list[1].Name[0]
		firstInt, _ := strconv.Atoi(string(nameFirstWordOfListZero))
		secondInt, _ := strconv.Atoi(string(nameFirstWordOfListOne))

		// Validate
		ts.GreaterOrEqual(secondInt, firstInt)
		ts.NoError(err)
	})

	ts.Run("Limit OFFSET", func() {
		pageModule := utils.NewPagination(1, 5)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, nil, common.ASC)

		// Validate
		ts.Equal(5, len(list))
		ts.NoError(err)
	})

	ts.Run("Limit OFFSET", func() {
		pageModule := utils.NewPagination(2, 20)
		list, err := ts.academyRepository.List(ts.ctx, pageModule, nil, nil, nil, nil, common.ASC)

		// Validate
		ts.Equal(0, len(list))
		ts.NoError(err)
	})

	ts.Run("Search By Yoga", func() {
		pageModule := utils.NewPagination(1, 10)
		yogaIDs := []int{2, 3}
		list, err := ts.academyRepository.List(ts.ctx, pageModule, &yogaIDs, nil, nil, nil, common.ASC)

		for _, v := range list {
			flag := false
			for _, id := range yogaIDs {
				if id == v.Edges.Yoga[0].ID {
					flag = true
					continue
				}
			}
			if !flag {
				ts.Fail("fail")
			}
		}
		// Validate
		ts.NoError(err)
	})
}

func TestAcademyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AcademyRepositoryTestSuite))
}

func (ts *AcademyRepositoryTestSuite) initForList() {
	bulk := make([]*ent.UserCreate, 40)

	for i := 0; i < 40; i++ {
		bulk[i] = ts.client.User.Create().
			SetEmail(fake.Email()).
			SetIsEmailVerified(true).
			SetNickname(fake.Animal()).
			SetPhoneNum(fake.Phone()).
			SetTermAgreeAt(time.Now())
	}
	err := ts.client.User.CreateBulk(bulk...).Exec(context.Background())
	ts.NoError(err)
	bulkForAcademy := make([]*ent.AcademyCreate, 20)

	ids, _ := ts.client.User.Query().Where(user.TypeIsNil()).Limit(20).IDs(context.Background())

	for i := 0; i < len(ids); i++ {
		sigunguID := 1
		if i%2 == 0 {
			sigunguID = 2
		}

		yogaid := i + 1
		if i > 7 {
			yogaid = 1
		}
		var yogaIDs []int

		yogaIDs = append(yogaIDs, yogaid)

		bulkForAcademy[i] = ts.client.Academy.Create().
			SetAddressDetail(fake.Address().Address).
			SetAddressRoad(fake.Address().Street).
			SetBusinessCode(fmt.Sprintf("%d", fake.Number(10000000, 9999999))).
			SetCallNumber(fake.Contact().Phone).
			SetName(fake.Animal()).SetUserID(ids[i]).SetSigunguID(sigunguID).AddYogaIDs(yogaIDs...)
	}
	err = ts.client.Academy.CreateBulk(bulkForAcademy...).Exec(context.Background())

	ts.NoError(err)
}

func (ts *AcademyRepositoryTestSuite) initData() {
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
