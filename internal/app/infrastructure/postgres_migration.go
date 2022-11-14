package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"

	"onthemat/internal/app/config"
	"onthemat/pkg/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func GeneratePostgresMigraion() {
	c := config.NewConfig()
	if err := c.Load("./configs"); err != nil {
		panic(err)
	}
	dir, err := atlas.NewLocalDir("internal/app/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                          // provide migration directory
		schema.WithMigrationMode(schema.ModeInspect), // provide migration mode
		schema.WithDialect(dialect.Postgres),         // Ent dialect to use
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
	}
	// Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.PostgreSQL.User, c.PostgreSQL.Password, c.PostgreSQL.Host, c.PostgreSQL.Port, c.PostgreSQL.Database)
	// url := "postgres://postgres:password@localhost:5432/test?sslmode=disable"
	err = migrate.NamedDiff(context.Background(), url, os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
