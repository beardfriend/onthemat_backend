package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"onthemat/internal/app/model"
	"onthemat/internal/app/utils"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/enttest"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserYogaRepository(t *testing.T) {
	c := utils.RepoTestInit(t)
	defer utils.RepoTestClose(t)
	time.Sleep(3 * time.Second)
	ctx := context.Background()
	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.User, c.PostgreSQL.Database, c.PostgreSQL.Password)
	client := enttest.Open(t, "postgres", url)
	defer client.Close()

	repo := NewUserYogaRepository(client)
	userRepo := NewUserRepository(client)

	tests := []utils.Tests{
		{
			Name: "userYogaRepository/createMany",

			Before: func(t *testing.T) {
				teacherType := model.TeacherType
				userRepo.Create(ctx, &ent.User{
					Type: &teacherType,
				})
			},

			Expect: func(t *testing.T) {
				t.Run("유저 아이디가 정상적으로 존재하는 경우", func(t *testing.T) {
					var yoga []*ent.UserYoga
					// 1개
					teacherType := model.TeacherType
					yoga = append(yoga, &ent.UserYoga{
						Name:     "하타",
						UserType: teacherType,
					})
					// 2개
					yoga = append(yoga, &ent.UserYoga{
						Name:     "인",
						UserType: teacherType,
					})
					us, err := repo.CreateMany(ctx, yoga, 1)
					assert.Equal(t, len(us), 2)
					assert.NoError(t, err)
				})

				t.Run("유저 아이디가 없는 경우", func(t *testing.T) {
					var yoga []*ent.UserYoga
					// 1개
					yoga = append(yoga, &ent.UserYoga{
						Name:     "하타",
						UserType: model.TeacherType,
					})

					_, err := repo.CreateMany(ctx, yoga, 2)
					assert.Equal(t, ent.IsConstraintError(err), true)
				})

				t.Run("빈 데이터만 입력했을 때", func(t *testing.T) {
					// 빈 배열을 넣어도 생성이 안되고 에러가 나지 않음.
					var yoga []*ent.UserYoga
					u, err := repo.CreateMany(ctx, yoga, 1)
					assert.Equal(t, len(u), 0)
					assert.NoError(t, err)
				})
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},

		{
			Name: "userYogaRepository/UpdateMany",

			Before: func(t *testing.T) {
				teacherType := model.TeacherType
				userRepo.Create(ctx, &ent.User{
					Type: &teacherType,
				})
			},

			Expect: func(t *testing.T) {
				t.Run("업데이트할 내역이 없는 경우", func(t *testing.T) {
				})

				t.Run("업데이트할 내역이 존재하는 경우", func(t *testing.T) {
				})

				t.Run("모든 값  삭제 요청이 들어온 경우.", func(t *testing.T) {
				})
			},

			After: func(t *testing.T) {
				utils.RepoTestTruncateTable(ctx, client)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Before(t)
			tt.Expect(t)
			tt.After(t)
		})
	}
}
