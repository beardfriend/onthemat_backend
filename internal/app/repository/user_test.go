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

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

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
			},
			After: func(t *testing.T) {
				utils.RepoTestRemoveTable(ctx, client)
			},
		},
		// Get
		{
			Name: "userRepository/Get",
			Before: func(t *testing.T) {
				_, err := userRepo.Create(ctx, &ent.User{})
				assert.NoError(t, err)
			},
			Expect: func(t *testing.T) {
			},
			After: func(t *testing.T) {
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

	utils.RepoTestClose(t)
}
