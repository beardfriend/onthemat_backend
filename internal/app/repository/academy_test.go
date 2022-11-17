package repository

import (
	"context"
	"testing"
	"time"

	"onthemat/internal/app/common"
	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"

	fakeit "github.com/brianvoe/gofakeit/v6"

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
		case "TestCreate":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.userId = u.ID

		case "TestGet":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.userId = u.ID
			err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
				Name:          "학원이름",
				BusinessCode:  "사업자번호",
				CallNumber:    "0226065418",
				AddressRoad:   "도로명주소",
				AddressSigun:  "서울시",
				AddressGu:     "강남구",
				AddressDong:   "논현동",
				AddressDetail: "상세주소",
				AddressX:      "x좌표",
				AddressY:      "y좌표",
			}, ts.userId)
			ts.NoError(err)

		case "TestUpdate":
			u, _ := ts.userRepo.Create(ts.ctx, &ent.User{})
			ts.userId = u.ID
			err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
				Name:          "학원이름",
				BusinessCode:  "사업자번호",
				CallNumber:    "0226065418",
				AddressRoad:   "도로명주소",
				AddressSigun:  "서울시",
				AddressGu:     "강남구",
				AddressDong:   "논현동",
				AddressDetail: "상세주소",
				AddressX:      "x좌표",
				AddressY:      "y좌표",
			}, ts.userId)
			ts.NoError(err)

		case "TestListAndTotal":

			userBulk := make([]*ent.UserCreate, 20)
			academyBulk := make([]*ent.AcademyCreate, 20)

			for i := 0; i < 20; i++ {
				userBulk[i] = ts.client.User.Create().
					SetNillableNickname(nil).
					SetNillableEmail(nil).
					SetNillablePassword(nil).
					SetNillableSocialKey(nil).
					SetNillableSocialName(nil).
					SetTermAgreeAt(time.Now()).
					SetNillablePhoneNum(nil)

				academyBulk[i] = ts.client.Academy.Create().
					SetAddressDetail(fakeit.Address().City).
					SetAddressDong(fakeit.Address().City).
					SetAddressGu(fakeit.Address().City).
					SetAddressRoad(fakeit.Address().City).
					SetAddressSigun(fakeit.Address().City).
					SetAddressX(fakeit.Address().City).
					SetAddressY(fakeit.Address().City).
					SetBusinessCode(fakeit.Address().City).
					SetCallNumber(fakeit.Phone()).
					SetName(fakeit.Address().City).SetUserID(i + 1)
			}

			err := ts.client.User.CreateBulk(userBulk...).Exec(ts.ctx)
			ts.NoError(err)
			err = ts.client.Academy.CreateBulk(academyBulk...).Exec(ts.ctx)
			ts.NoError(err)

		}
	}
}

func (ts *AcademyRepositoryTestSuite) TestCreate() {
	ts.Run("성공", func() {
		err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
			Name:          "학원이름",
			BusinessCode:  "사업자번호",
			CallNumber:    "0226065418",
			AddressRoad:   "도로명주소",
			AddressSigun:  "서울시",
			AddressGu:     "강남구",
			AddressDong:   "논현동",
			AddressDetail: "상세주소",
			AddressX:      "x좌표",
			AddressY:      "y좌표",
		}, ts.userId)
		ts.NoError(err)
	})

	ts.Run("없는 유저일 경우.", func() {
		err := ts.academyRepository.Create(ts.ctx, &ent.Academy{
			Name:          "학원이름",
			BusinessCode:  "사업자번호",
			CallNumber:    "0226065418",
			AddressRoad:   "도로명주소",
			AddressSigun:  "서울시",
			AddressGu:     "강남구",
			AddressDong:   "논현동",
			AddressDetail: "상세주소",
			AddressX:      "x좌표",
			AddressY:      "y좌표",
		}, 2)
		ts.Equal(true, ent.IsConstraintError(err))
	})
}

func (ts *AcademyRepositoryTestSuite) TestGet() {
	ts.Run("성공", func() {
		academy, err := ts.academyRepository.Get(ts.ctx, ts.userId)
		ts.Equal(academy.Name, "학원이름")
		ts.NoError(err)
	})
}

func (ts *AcademyRepositoryTestSuite) TestUpdate() {
	ts.Run("성공", func() {
		err := ts.academyRepository.Update(ts.ctx, &ent.Academy{
			Name:          "이름변경",
			BusinessCode:  "사업자번호",
			CallNumber:    "0226065418",
			AddressRoad:   "도로명주소",
			AddressSigun:  "서울시",
			AddressGu:     "강남구",
			AddressDong:   "논현동",
			AddressDetail: "상세주소",
			AddressX:      "x좌표",
			AddressY:      "y좌표",
		}, ts.userId)
		ts.NoError(err)

		a, _ := ts.academyRepository.Get(ts.ctx, ts.userId)
		ts.Equal(a.Name, "이름변경")
	})

	ts.Run("없는 유저", func() {
		// 쿼리 시 where로 찾으면 notFound error가 나오지 않고, UpdateOneId로 찾으면 notfound에러가 떨어짐
		err := ts.academyRepository.Update(ts.ctx, &ent.Academy{
			Name:          "이름변경",
			BusinessCode:  "사업자번호",
			CallNumber:    "0226065418",
			AddressRoad:   "도로명주소",
			AddressSigun:  "서울시",
			AddressGu:     "강남구",
			AddressDong:   "논현동",
			AddressDetail: "상세주소",
			AddressX:      "x좌표",
			AddressY:      "y좌표",
		}, 2)
		ts.Equal(true, ent.IsNotFound(err))
	})
}

func (ts *AcademyRepositoryTestSuite) TestListAndTotal() {
	ts.Run("검색할 수 없는 컬럼", func() {
		searchKey := "HaventExisitColumn"
		searchValue := "nasd"
		_, err := ts.academyRepository.List(ts.ctx, &common.ListParams{
			PageNo:      1,
			PageSize:    10,
			SearchKey:   &searchKey,
			SearchValue: &searchValue,
			OrderType:   nil,
			OrderCol:    nil,
		})
		ts.Error(err, ErrSearchColumnInvalid)
	})

	ts.Run("페이지 사이즈 11개", func() {
		a, err := ts.academyRepository.List(ts.ctx, &common.ListParams{
			PageNo:      1,
			PageSize:    11,
			SearchKey:   nil,
			SearchValue: nil,
			OrderType:   nil,
			OrderCol:    nil,
		})
		ts.NoError(err)
		ts.Equal(len(a), 11)
	})

	ts.Run("페이지 넘버 2", func() {
		a, err := ts.academyRepository.List(ts.ctx, &common.ListParams{
			PageNo:      2,
			PageSize:    15,
			SearchKey:   nil,
			SearchValue: nil,
			OrderType:   nil,
			OrderCol:    nil,
		})
		ts.NoError(err)
		ts.Equal(len(a), 5)
	})

	ts.Run("토탈", func() {
		a, err := ts.academyRepository.Total(ts.ctx, &common.TotalParams{
			SearchKey:   nil,
			SearchValue: nil,
		})
		ts.NoError(err)
		ts.Equal(a, 20)
	})
}

func TestAcademyRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AcademyRepositoryTestSuite))
}
