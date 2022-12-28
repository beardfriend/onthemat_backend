package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/model"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/yoga"

	fake "github.com/brianvoe/gofakeit/v6"
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
}

// 모든 테스트 시작 전 1회
func (ts *YogaRepositoryTestSuite) SetupSuite() {
	t := ts.T()
	ts.ctx = context.Background()
	c := utils.GetTestConfig(t)
	ts.config = c
	utils.RepoTestClose(t)
	time.Sleep(1 * time.Second)
	ts.config = utils.RepoTestInit(t)
	time.Sleep(3 * time.Second)

	// 포스트그레스 연결
	ts.client = infrastructure.NewPostgresDB(ts.config)
	elastic := infrastructure.NewElasticSearch(ts.config, "../../../configs/elastic.crt")

	// 모듈 연결
	ts.yogaRepo = NewYogaRepository(ts.client, elastic)
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

func (ts *YogaRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "YogaRepositoryTestSuite" {
		switch testName {
		case "TestCreateGroup":

		case "TestUpdateGroup":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestDeleteGroups":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "하타",
				CategoryEng: "hata",
				Description: utils.String("하타 요가입니다."),
			})
			ts.NoError(err)

			_, err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "인",
				CategoryEng: "in",
				Description: utils.String("인 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "플라잉",
				CategoryEng: "flying",
				Description: utils.String("플라잉 요가입니다."),
			})
			ts.NoError(err)

			_, err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				YogaGroupID: 4,
				NameKor:     "플라잉 키즈",
			})
			ts.NoError(err)

		case "TestPatchGroup":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestGroupTotal":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "인",
				CategoryEng: "in",
				Description: utils.String("인 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "플라잉",
				CategoryEng: "flying",
				Description: utils.String("플라잉 요가입니다."),
			})
			ts.NoError(err)

		case "TestGroupList":
			bulk := make([]*ent.YogaGroupCreate, 20)
			for i := 0; i < 20; i++ {
				category := fmt.Sprintf("%s%d", fake.BeerName(), i)
				bulk[i] = ts.client.YogaGroup.
					Create().
					SetCategory(category).
					SetCategoryEng(fake.BeerHop())
			}
			_, err := ts.client.YogaGroup.CreateBulk(bulk...).Save(ts.ctx)
			ts.NoError(err)

		case "TestCreate":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestUpdate":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				YogaGroupID: 1,
			})

			ts.NoError(err)

		case "TestPatch":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				Level:       utils.Int(1),
				Description: utils.String("description"),
				YogaGroupID: 1,
			})
			ts.NoError(err)
		case "TestDelete":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			_, err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				Level:       utils.Int(1),
				Description: utils.String("description"),
				YogaGroupID: 1,
			})
			ts.NoError(err)

		case "TestList":
			_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			bulk := make([]*ent.YogaCreate, 20)
			for i := 0; i < 20; i++ {
				bulk[i] = ts.client.Yoga.Create().
					SetNameKor(fake.Animal()).SetYogaGroupID(1)
			}
			_, err = ts.client.Yoga.CreateBulk(bulk...).Save(ts.ctx)
			ts.NoError(err)

		case "TestElasticList":
			err := ts.yogaRepo.ElasticCreate(ts.ctx, &model.ElasticYoga{
				Id:   1,
				Name: "아쉬탕가",
			})
			ts.NoError(err)

			err = ts.yogaRepo.ElasticCreate(ts.ctx, &model.ElasticYoga{
				Id:   2,
				Name: "아쉬탕가 레드",
			})
			ts.NoError(err)

			err = ts.yogaRepo.ElasticCreate(ts.ctx, &model.ElasticYoga{
				Id:   3,
				Name: "아쉬탕가 프리이머리",
			})
			ts.NoError(err)

		case "TestElasticUpdate":
			err := ts.yogaRepo.ElasticCreate(ts.ctx, &model.ElasticYoga{
				Id:   1,
				Name: "아쉬탕가",
			})
			ts.NoError(err)
		}
	}
}

func (ts *YogaRepositoryTestSuite) AfterTest(suiteName, testName string) {
	if suiteName == "YogaRepositoryTestSuite" {
		if testName == "TestElasticList" {
			rowAffected, err := ts.yogaRepo.ElasitcDelete(ts.ctx, []int{1, 2, 3})
			ts.NoError(err)
			ts.Equal(3, rowAffected)
		}
	}
}

// TEST GROUP
func (ts *YogaRepositoryTestSuite) TestCreateGroup() {
	ts.Run("success", func() {
		_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
			Category:    "아쉬탕가",
			CategoryEng: "ashtanga",
			Description: utils.String("아쉬탕가 요가입니다."),
		})
		ts.NoError(err)

		res, _ := ts.client.YogaGroup.Get(ts.ctx, 1)

		ts.Equal(res.Category, "아쉬탕가")
	})
	ts.Run("null값으로 들어가는지.", func() {
		_, err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
			Category:    "하타",
			CategoryEng: "Hata",
			Description: nil,
		})
		ts.NoError(err)

		res, _ := ts.client.YogaGroup.Get(ts.ctx, 2)
		ts.Nil(res.Description)
	})
}

func (ts *YogaRepositoryTestSuite) TestUpdateGroup() {
	ts.Run("success", func() {
		err := ts.yogaRepo.UpdateGroup(ts.ctx, &ent.YogaGroup{
			ID:       1,
			Category: "아쉬람",
		})
		ts.NoError(err)

		res, _ := ts.client.YogaGroup.Get(ts.ctx, 1)
		ts.Equal(res.Category, "아쉬람")
		ts.Equal(res.CategoryEng, "")
		ts.Nil(res.Description)
	})
}

func (ts *YogaRepositoryTestSuite) TestPatchGroup() {
	ts.Run("success", func() {
		err := ts.yogaRepo.PatchGroup(ts.ctx, &request.YogaGroupPatchBody{
			Category: utils.String("하하요가"),
		}, 1)
		ts.NoError(err)

		res, _ := ts.client.YogaGroup.Get(ts.ctx, 1)
		ts.Equal("하하요가", res.Category)
		ts.Equal("ashtanga", res.CategoryEng)
		ts.Equal("아쉬탕가 요가입니다.", *res.Description)
	})
}

func (ts *YogaRepositoryTestSuite) TestDeleteGroups() {
	ts.Run("success", func() {
		rowAffected, err := ts.yogaRepo.DeleteGroups(ts.ctx, []int{1, 2})
		ts.NoError(err)
		ts.Equal(2, rowAffected)
	})

	ts.Run("yoga Group밑에 요가가 몇가지 등록되어 있다면", func() {
		rowAffected, err := ts.yogaRepo.DeleteGroups(ts.ctx, []int{3, 4})
		ts.NoError(err)
		ts.Equal(2, rowAffected)

		res, err := ts.client.Yoga.Get(ts.ctx, 1)
		ts.NoError(err)
		ts.Equal(0, res.YogaGroupID)
	})
}

func (ts *YogaRepositoryTestSuite) TestGroupTotal() {
	ts.Run("success", func() {
		count, err := ts.yogaRepo.GroupTotal(ts.ctx, nil)
		ts.NoError(err)
		ts.Equal(2, count)
	})

	ts.Run("category", func() {
		count, err := ts.yogaRepo.GroupTotal(ts.ctx, utils.String("인"))
		ts.NoError(err)
		ts.Equal(1, count)
	})
}

func (ts *YogaRepositoryTestSuite) TestGroupList() {
	ts.Run("success", func() {
		pgModule := utils.NewPagination(1, 10)
		result, err := ts.yogaRepo.GroupList(ts.ctx, pgModule, nil, common.ASC)
		ts.NoError(err)
		ts.Equal(10, len(result))
		ts.Greater(result[1].ID, result[0].ID)
	})

	ts.Run("query", func() {
		pgModule := utils.NewPagination(1, 10)
		result, err := ts.yogaRepo.GroupList(ts.ctx, pgModule, utils.String("0"), common.DESC)
		ts.NoError(err)
		ts.GreaterOrEqual(len(result), 1)
	})
}

func (ts *YogaRepositoryTestSuite) TestCreate() {
	ts.Run("success", func() {
		_, err := ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
			NameKor:     "아쉬탕가 프라이머리",
			YogaGroupID: 1,
		})
		ts.NoError(err)

		res, _ := ts.client.Yoga.Get(ts.ctx, 1)

		ts.Equal("아쉬탕가 프라이머리", res.NameKor)
		ts.Equal(1, res.YogaGroupID)
		ts.Nil(res.Level)
		ts.Nil(res.NameEng)
	})
}

func (ts *YogaRepositoryTestSuite) TestUpdate() {
	ts.Run("success", func() {
		ts.yogaRepo.Update(ts.ctx, &ent.Yoga{
			ID:          1,
			Level:       utils.Int(1),
			YogaGroupID: 1,
		})

		res, _ := ts.client.Yoga.Get(ts.ctx, 1)

		ts.Equal("", res.NameKor)
		ts.Equal(1, res.YogaGroupID)
		ts.Equal(1, *res.Level)
		ts.Nil(res.NameEng)
	})

	ts.Run("없는 아이디", func() {
		err := ts.yogaRepo.Update(ts.ctx, &ent.Yoga{
			ID:          3,
			Level:       utils.Int(1),
			YogaGroupID: 1,
		})
		ts.NoError(err)
	})
}

func (ts *YogaRepositoryTestSuite) TestPatch() {
	ts.Run("success", func() {
		ts.yogaRepo.Patch(ts.ctx, &request.YogaPatchBody{
			Level: utils.Int(5),
		}, 1)

		res, _ := ts.client.Yoga.Get(ts.ctx, 1)

		ts.Equal("아쉬탕가 프라이머리", res.NameKor)
		ts.Equal(1, res.YogaGroupID)
		ts.Equal(5, *res.Level)
		ts.Equal("ashtanga primary", *res.NameEng)
		ts.Equal("description", *res.Description)
	})
}

func (ts *YogaRepositoryTestSuite) TestDelete() {
	ts.Run("success", func() {
		err := ts.yogaRepo.Delete(ts.ctx, 1)
		ts.NoError(err)

		_, err = ts.client.Yoga.Query().Where(yoga.IDEQ(1)).Only(ts.ctx)
		ent.IsNotFound(err)

		res, err := ts.client.YogaGroup.Get(ts.ctx, 1)
		ts.NoError(err)
		ts.Equal(1, res.ID)
	})
}

func (ts *YogaRepositoryTestSuite) TestList() {
	ts.Run("success", func() {
		/*
			SELECT DISTINCT "yoga"."id", "yoga"."created_at", "yoga"."updated_at", "yoga"."yoga_group_id", "yoga"."name_kor", "yoga"."name_eng", "yoga"."level", "yoga"."description"
			FROM "yoga"
			JOIN (SELECT "yoga_group"."id" FROM "yoga_group" WHERE "yoga_group"."id" = $1) AS "t1"
			ON "yoga"."yoga_group_id" = "t1"."id"
			ORDER BY "yoga"."id" DESC args=[1]
		*/
		result, err := ts.yogaRepo.List(ts.ctx, 1)
		ts.NoError(err)
		ts.Equal(20, len(result))
	})
}

func (ts *YogaRepositoryTestSuite) TestElasticList() {
	ts.Run("success", func() {
		d, err := ts.yogaRepo.ElasticList(ts.ctx, "아쉬탕가 레")
		ts.NoError(err)
		fmt.Println(d)
	})
}

func (ts *YogaRepositoryTestSuite) TestElasticUpdate() {
	ts.Run("success", func() {
		rawAffected, err := ts.yogaRepo.ElasticUpdate(ts.ctx, &model.ElasticYoga{
			Id:   1,
			Name: "아헹가",
		})
		ts.NoError(err)
		ts.Equal(1, rawAffected)
		time.Sleep(2 * time.Second)
		d, err := ts.yogaRepo.ElasticList(ts.ctx, "아헹가")
		ts.NoError(err)
		fmt.Println(d)
	})
}

func TestYogaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(YogaRepositoryTestSuite))
}
