package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/repository"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/enttest"
	"onthemat/pkg/ent/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	utils.RepoTestClose(t)
}

func TestUserRepository(t *testing.T) {
	// init
	c := utils.RepoTestInit(t)
	time.Sleep(1 * time.Second)
	ctx := context.Background()
	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.User, c.PostgreSQL.Database, c.PostgreSQL.Password)
	client := enttest.Open(t, "postgres", url)
	defer client.Close()
	userRepo := repository.NewUserRepository(client)

	tests := []utils.Tests{
		// Create
		{
			Name: "userRepository/Create",

			Before: func(t *testing.T) {
			},

			Expect: func(t *testing.T) {
				t.Run("ID만 있는 경우", func(t *testing.T) {
					u, err := userRepo.Create(ctx, &ent.User{})
					assert.NoError(t, err)
					assert.Empty(t, u.Email)
					assert.Empty(t, u.Nickname)
					assert.Empty(t, u.Password)
					assert.Empty(t, u.PhoneNum)
					assert.Empty(t, u.RemovedAt)
					assert.Empty(t, u.SocialKey)
					assert.Empty(t, u.SocialName)
					assert.Equal(t, u.ID, 1)
				})

				t.Run("적당한 데이터", func(t *testing.T) {
					data := struct {
						Email     string
						Password  string
						Nickname  string
						PhoneNum  string
						TermAgree time.Time
					}{
						Email:     "asd@naver.com",
						Password:  "password",
						Nickname:  "nickname",
						PhoneNum:  "01043226633",
						TermAgree: time.Now(),
					}

					_, err := userRepo.Create(ctx, &ent.User{
						Email:       &data.Email,
						Password:    &data.Password,
						Nickname:    &data.Nickname,
						PhoneNum:    &data.PhoneNum,
						TermAgreeAt: &data.TermAgree,
					})

					assert.NoError(t, err)
				})

				t.Run("SocialKey 중복됐을 때", func(t *testing.T) {
					data := struct {
						SocialName user.SocialName
						SocialKey  string
					}{
						SocialName: user.SocialNameKakao,
						SocialKey:  "asdadsads",
					}

					_, err := userRepo.Create(ctx, &ent.User{
						SocialName: &data.SocialName,
						SocialKey:  &data.SocialKey,
					})

					assert.NoError(t, err)

					// duplicated key
					_, err = userRepo.Create(ctx, &ent.User{
						SocialName: &data.SocialName,
						SocialKey:  &data.SocialKey,
					})

					assert.Equal(t, ent.IsConstraintError(err), true)
				})
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},
		// Update
		{
			Name: "userRepository/Update",

			Before: func(t *testing.T) {
				email := "asd@naver.com"
				password := "password"
				nickname := "nick"
				phoneNum := "01043226633"
				termAgreeAt := time.Now()
				userRepo.Create(ctx, &ent.User{
					Email:       &email,
					Password:    &password,
					Nickname:    &nickname,
					PhoneNum:    &phoneNum,
					TermAgreeAt: &termAgreeAt,
				})
			},

			Expect: func(t *testing.T) {
				email := "das@naver.com"
				nickname := "kick"
				phoneNum := "01064135418"
				uu, err := userRepo.Update(ctx, &ent.User{
					ID:       1,
					Email:    &email,
					Nickname: &nickname,
					PhoneNum: &phoneNum,
				})

				assert.NoError(t, err)
				assert.Equal(t, *uu.Email, email)
				assert.Equal(t, *uu.Nickname, nickname)
				assert.Equal(t, *uu.PhoneNum, phoneNum)
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},

		// Get
		{
			Name: "userRepository/Get",

			Before: func(t *testing.T) {
				email := "asd@gmail.com"
				_, err := userRepo.Create(ctx, &ent.User{
					Email: &email,
				})
				assert.NoError(t, err)
			},

			Expect: func(t *testing.T) {
				u, err := userRepo.Get(ctx, 1)
				assert.NoError(t, err)
				assert.Equal(t, *u.Email, "asd@gmail.com")
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},

		// FindByEmail
		{
			Name: "userRepository/FindByEmail",

			Before: func(t *testing.T) {
				email := "asd@gmail.com"
				_, err := userRepo.Create(ctx, &ent.User{
					Email: &email,
				})
				assert.NoError(t, err)
			},

			Expect: func(t *testing.T) {
				t.Run("success", func(t *testing.T) {
					isExist, err := userRepo.FindByEmail(ctx, "asd@gmail.com")
					assert.NoError(t, err)
					assert.Equal(t, isExist, true)
				})

				t.Run("no", func(t *testing.T) {
					isExist, err := userRepo.FindByEmail(ctx, "noasd@gmail.com")
					assert.NoError(t, err)
					assert.Equal(t, isExist, false)
				})
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},

		// GetByEmailPassword
		{
			Name: "userRepository/GetByEmailPassword",

			Before: func(t *testing.T) {
				email := "asd@gmail.com"
				password := "password"
				_, err := userRepo.Create(ctx, &ent.User{
					Email:    &email,
					Password: &password,
				})
				assert.NoError(t, err)
			},

			Expect: func(t *testing.T) {
				email := "asd@gmail.com"
				password := "password"

				t.Run("success", func(t *testing.T) {
					u, err := userRepo.GetByEmailPassword(ctx, &ent.User{
						Email:    &email,
						Password: &password,
					})
					assert.NoError(t, err)
					assert.Equal(t, *u.Email, email)
				})
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},

		//
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Before(t)
			tt.Expect(t)
			tt.After(t)
		})
	}

	utils.RepoTestClose(t)
}
