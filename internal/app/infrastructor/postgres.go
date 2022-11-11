package infrastructor

import (
	"context"
	"log"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/migrate"

	_ "github.com/lib/pq"
)

func NewPostgresDB() *ent.Client {
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=db password=password sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

func ClosePostgres(client *ent.Client) error {
	return client.Close()
}
