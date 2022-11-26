package repository

import (
	"context"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/yoga"

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

	// 포스트그레스 연결
	ts.client = infrastructure.NewPostgresDB(ts.config)

	// 모듈 연결
	ts.yogaRepo = NewYogaRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// 모든 테스트 종료 후 1회
func (ts *YogaRepositoryTestSuite) TearDownSuite() {
	ts.client.Close()
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
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestDeleteGroups":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "하타",
				CategoryEng: "hata",
				Description: utils.String("하타 요가입니다."),
			})
			ts.NoError(err)

			err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "인",
				CategoryEng: "in",
				Description: utils.String("인 요가입니다."),
			})
			ts.NoError(err)
			err = ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "플라잉",
				CategoryEng: "flying",
				Description: utils.String("플라잉 요가입니다."),
			})
			ts.NoError(err)

			err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				YogaGroupID: 4,
				NameKor:     "플라잉 키즈",
			})
			ts.NoError(err)

		case "TestPatchGroup":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestCreate":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)

		case "TestUpdate":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				YogaGroupID: 1,
			})

			ts.NoError(err)

		case "TestPatch":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				Level:       utils.Int(1),
				Description: utils.String("description"),
				YogaGroupID: 1,
			})
			ts.NoError(err)
		case "TestDelete":
			err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
				Category:    "아쉬탕가",
				CategoryEng: "ashtanga",
				Description: utils.String("아쉬탕가 요가입니다."),
			})
			ts.NoError(err)
			err = ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
				NameKor:     "아쉬탕가 프라이머리",
				NameEng:     utils.String("ashtanga primary"),
				Level:       utils.Int(1),
				Description: utils.String("description"),
				YogaGroupID: 1,
			})
			ts.NoError(err)

		}
	}
}

// TEST GROUP
func (ts *YogaRepositoryTestSuite) TestCreateGroup() {
	ts.Run("success", func() {
		err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
			Category:    "아쉬탕가",
			CategoryEng: "ashtanga",
			Description: utils.String("아쉬탕가 요가입니다."),
		})
		ts.NoError(err)

		res, _ := ts.client.YogaGroup.Get(ts.ctx, 1)

		ts.Equal(res.Category, "아쉬탕가")
	})
	ts.Run("null값으로 들어가는지.", func() {
		err := ts.yogaRepo.CreateGroup(ts.ctx, &ent.YogaGroup{
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

func (ts *YogaRepositoryTestSuite) TestCreate() {
	ts.Run("success", func() {
		err := ts.yogaRepo.Create(ts.ctx, &ent.Yoga{
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

func TestYogaRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(YogaRepositoryTestSuite))
}
