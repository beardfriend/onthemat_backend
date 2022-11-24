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

type UserRepositoryTestSuite struct {
	suite.Suite
	config   *config.Config
	client   *ent.Client
	userRepo UserRepository
	yogaRepo YogaRepository
	ctx      context.Context
	// Data
	testUpdateData struct {
		id          int
		email       string
		password    string
		nickname    string
		phoneNum    string
		termAgreeAt time.Time
	}

	testGetData struct {
		id    int
		email string
	}

	testFindByEmailData struct {
		id    int
		email string
	}

	testGetByEmailPassword [3]struct {
		id           int
		email        string
		password     string
		tempPassword string
		nickname     string
	}

	testAddYogaData struct {
		id    int
		email string
	}

	// Flag For BeforeRun
	createSocialKey bool
}

// ------------------- Running Before Test Start Once   -------------------

func (ts *UserRepositoryTestSuite) SetupSuite() {
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
	ts.yogaRepo = NewYogaRepository(ts.client)
	ts.userRepo = NewUserRepository(ts.client)
}

// ------------------- Running After Every Test Finish  -------------------

func (ts *UserRepositoryTestSuite) TearDownSuite() {
	// 도커 내리기
	ts.client.Close()
	utils.RepoTestClose(nil)
}

// ------------------- Running After Each Test Finish  -------------------

func (ts *UserRepositoryTestSuite) TearDownTest() {
	// 테이블 Truncate
	utils.RepoTestTruncateTable(context.Background(), ts.client)
}

// ------------------- SetUp Data Before Each Test Start-------------------

func (ts *UserRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	if suiteName == "UserRepositoryTestSuite" {
		switch testName {
		case "TestUpdate":
			ts.testUpdateData.email = "asd@naver.com"
			ts.testUpdateData.password = "password"
			ts.testUpdateData.nickname = "nick"
			ts.testUpdateData.phoneNum = "01043226633"
			ts.testUpdateData.termAgreeAt = time.Now()

			u, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email:       &ts.testUpdateData.email,
				Password:    &ts.testUpdateData.password,
				Nickname:    &ts.testUpdateData.nickname,
				PhoneNum:    &ts.testUpdateData.phoneNum,
				TermAgreeAt: ts.testUpdateData.termAgreeAt,
			})
			ts.NoError(err)
			ts.testUpdateData.id = u.ID
		case "TestFindByEmail":
			ts.testFindByEmailData.email = "asd@gmail.com"
			user, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email: &ts.testFindByEmailData.email,
			})
			ts.NoError(err)

			ts.testFindByEmailData.id = user.ID

		case "TestGet":
			ts.testGetData.email = "asd@gmail.com"
			user, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email: &ts.testFindByEmailData.email,
			})
			ts.NoError(err)

			ts.testGetData.id = user.ID

		case "TestGetByEmaillPassword":
			// 첫 번째 유저

			ts.testGetByEmailPassword[0].email = "asd@gmail.com"
			ts.testGetByEmailPassword[0].password = "password"
			user1, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email:    &ts.testGetByEmailPassword[0].email,
				Password: &ts.testGetByEmailPassword[0].password,
			})
			ts.NoError(err)
			ts.testGetByEmailPassword[0].id = user1.ID

			// 임시비밀번호 유저

			ts.testGetByEmailPassword[1].email = "asd2@gmail.com"
			ts.testGetByEmailPassword[1].password = "password"
			ts.testGetByEmailPassword[1].tempPassword = "tempPassword"
			ts.testGetByEmailPassword[1].nickname = "nickname"
			user2, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email:    &ts.testGetByEmailPassword[1].email,
				Password: &ts.testGetByEmailPassword[1].password,
				Nickname: &ts.testGetByEmailPassword[1].nickname,
			})
			ts.NoError(err)

			err = ts.userRepo.UpdateTempPassword(ts.ctx, &ent.User{
				Email:        &ts.testGetByEmailPassword[1].email,
				TempPassword: &ts.testGetByEmailPassword[1].tempPassword,
			})
			ts.NoError(err)
			ts.testGetByEmailPassword[1].id = user2.ID

		case "TestAddYoga":
			ts.testAddYogaData.email = "asd@gmail.com"
			user, err := ts.userRepo.Create(ts.ctx, &ent.User{
				Email: &ts.testFindByEmailData.email,
			})
			ts.NoError(err)

			ts.testAddYogaData.id = user.ID
		}
	}
}

func (ts *UserRepositoryTestSuite) TestCreate() {
	ts.Run("ID만 존재하는 경우", func() {
		u, err := ts.userRepo.Create(ts.ctx, &ent.User{})
		ts.NoError(err)
		ts.Empty(u.Email)
		ts.Empty(u.Nickname)
		ts.Empty(u.Password)
		ts.Empty(u.PhoneNum)
		ts.Empty(u.SocialKey)
		ts.Empty(u.SocialName)
		ts.Equal(u.ID, 1)
	})

	ts.Run("일반 회원가입 시 들어오는 정보들", func() {
		email := "asd@naver.com"
		password := "password"
		nickname := "nickname"
		phoneNum := "01043226633"
		TermAgree := time.Now()

		_, err := ts.userRepo.Create(ts.ctx, &ent.User{
			Email:       &email,
			Password:    &password,
			Nickname:    &nickname,
			PhoneNum:    &phoneNum,
			TermAgreeAt: TermAgree,
		})
		ts.NoError(err)
	})

	ts.Run("SocialKey 중복됐을 때", func() {
		ts.createSocialKey = true
		socialName := model.KakaoSocialType
		socialKey := "asdasdasd"

		_, err := ts.userRepo.Create(ts.ctx, &ent.User{
			SocialName: &socialName,
			SocialKey:  &socialKey,
		})

		ts.NoError(err)

		// duplicated key
		_, err = ts.userRepo.Create(ts.ctx, &ent.User{
			SocialName: &socialName,
			SocialKey:  &socialKey,
		})

		ts.Equal(ent.IsConstraintError(err), true)
	})

	ts.Run("이메일 중복됐을 때", func() {
		email := "asd123@naver.com"

		_, err := ts.userRepo.Create(ts.ctx, &ent.User{
			Email: &email,
		})

		ts.NoError(err)

		// duplicated key
		_, err = ts.userRepo.Create(ts.ctx, &ent.User{
			Email: &email,
		})

		ts.Equal(ent.IsConstraintError(err), true)
	})
}

func (ts *UserRepositoryTestSuite) TestUpdate() {
	ts.Run("성공", func() {
		nickname := "kick"
		phoneNum := "01064135418"

		user, err := ts.userRepo.Update(ts.ctx, &ent.User{
			ID:       ts.testUpdateData.id,
			Nickname: &nickname,
			PhoneNum: &phoneNum,
		})

		ts.NoError(err)
		ts.Equal(*user.Nickname, nickname)
		ts.Equal(*user.PhoneNum, phoneNum)
	})
}

func (ts *UserRepositoryTestSuite) TestGet() {
	ts.Run("성공", func() {
		user, err := ts.userRepo.Get(ts.ctx, ts.testGetData.id)
		ts.NoError(err)
		ts.Equal(*user.Email, ts.testGetData.email)
	})

	ts.Run("Not Found Error", func() {
		_, err := ts.userRepo.Get(ts.ctx, 2)
		ts.Equal(ent.IsNotFound(err), true)
	})
}

func (ts *UserRepositoryTestSuite) TestFindByEmail() {
	ts.Run("존재하는 이메일", func() {
		isExist, err := ts.userRepo.FindByEmail(ts.ctx, ts.testFindByEmailData.email)
		ts.NoError(err)
		ts.Equal(isExist, true)
	})

	ts.Run("존재하지 않는 이메일", func() {
		isExist, err := ts.userRepo.FindByEmail(ts.ctx, "nirvana@buddha.com")
		ts.NoError(err)
		ts.Equal(isExist, false)
	})
}

func (ts *UserRepositoryTestSuite) TestGetByEmaillPassword() {
	ts.Run("이메일 비밀번호 모두 일치할 때", func() {
		user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
			Email:    &ts.testGetByEmailPassword[0].email,
			Password: &ts.testGetByEmailPassword[0].password,
		})
		ts.NoError(err)
		ts.Equal(*user.Email, ts.testGetByEmailPassword[0].email)
	})

	ts.Run("임시 비밀번호가 있는 경우", func() {
		ts.Run("임시 비밀번호로 조회", func() {
			user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
				Email:    &ts.testGetByEmailPassword[1].email,
				Password: &ts.testGetByEmailPassword[1].tempPassword,
			})
			ts.NoError(err)
			ts.Equal(*user.Email, ts.testGetByEmailPassword[1].email)
		})

		ts.Run("일반 비밀번호로 조회", func() {
			user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
				Email:    &ts.testGetByEmailPassword[1].email,
				Password: &ts.testGetByEmailPassword[1].password,
			})
			ts.NoError(err)
			ts.Equal(*user.Email, ts.testGetByEmailPassword[1].email)
		})
	})

	ts.Run("이메일 혹은 비밀번호가 일치하지 않을 때", func() {
		noEmail := "noEmail@no.com"
		inCorrectPassword := "incorrectPassword"

		ts.Run("이메일 불일치", func() {
			user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
				Email:    &noEmail,
				Password: &ts.testGetByEmailPassword[0].password,
			})
			ts.Nil(user)
			ts.Equal(ent.IsNotFound(err), true)
		})

		ts.Run("패스워드 불일치", func() {
			user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
				Email:    &ts.testGetByEmailPassword[0].email,
				Password: &inCorrectPassword,
			})
			ts.Nil(user)
			ts.Equal(ent.IsNotFound(err), true)
		})

		ts.Run("둘 다 불일치", func() {
			user, err := ts.userRepo.GetByEmailPassword(ts.ctx, &ent.User{
				Email:    &noEmail,
				Password: &inCorrectPassword,
			})
			ts.Nil(user)
			ts.Equal(ent.IsNotFound(err), true)
		})
	})
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
