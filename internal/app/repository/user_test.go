package repository_test

import (
	"context"
	"fmt"
	"testing"

	"onthemat/internal/app/repository"
	"onthemat/pkg/ent"
	"onthemat/pkg/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name   string
		before func(*testing.T, *ent.Client)
		expect func(*testing.T, repository.UserRepository)
		after  func(*ent.Client)
	}{
		{
			name: "create",

			expect: func(t *testing.T, userRepo repository.UserRepository) {
				_, err := userRepo.Create(ctx, &ent.User{})
				assert.NoError(t, err)
			},
			after: func(c *ent.Client) {
				c.ExecContext(ctx, `
				DO $$ DECLARE
				r RECORD;
				BEGIN
				FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
					EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
				END LOOP;
			END $$;
				`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := enttest.Open(t, "postgres", "host=localhost port=5432 user=postgres dbname=test password=password sslmode=disable")
			defer cli.Close()
			userRepo := repository.NewUserRepository(cli)

			tt.expect(t, userRepo)
			tt.after(cli)
		})
	}
}

func TestRepo_FindByEmail(t *testing.T) {
	ctx := context.Background()
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	assert.NoError(t, err)
	defer client.Close()
	err = client.Schema.Create(context.Background())
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(client)
	nick := "asd"
	user, err := userRepo.Create(ctx, &ent.User{
		Nickname: &nick,
	})
	assert.NoError(t, err)

	fmt.Println(user.Nickname)
}

func TestRepo_FindByEmailPg(t *testing.T) {
	ctx := context.Background()
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=test password=password sslmode=disable")
	assert.NoError(t, err)
	// Run the auto migration tool.
	err = client.Schema.Create(context.Background())
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(client)

	nick := "asd"
	user, err := userRepo.Create(ctx, &ent.User{
		Nickname: &nick,
	})
	assert.NoError(t, err)

	fmt.Println(user.Nickname)
}

func TestRepo_Get(t *testing.T) {
	ctx := context.Background()
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=test password=password sslmode=disable")
	assert.NoError(t, err)
	// Run the auto migration tool.
	err = client.Schema.Create(context.Background())
	assert.NoError(t, err)

	userRepo := repository.NewUserRepository(client)
	u, err := userRepo.Get(ctx, 2)

	fmt.Println(u)
}
