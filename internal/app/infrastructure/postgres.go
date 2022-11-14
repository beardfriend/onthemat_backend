package infrastructure

import (
	"context"
	"fmt"
	"log"

	"onthemat/internal/app/config"
	"onthemat/pkg/ent"

	_ "github.com/lib/pq"
)

func NewPostgresDB(c *config.Config) *ent.Client {
	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.User, c.PostgreSQL.Database, c.PostgreSQL.Password)
	client, err := ent.Open("postgres", url)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return client
}

func ClosePostgres(client *ent.Client) error {
	return client.Close()
}
