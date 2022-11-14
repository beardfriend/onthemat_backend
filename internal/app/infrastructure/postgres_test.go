package infrastructure

import (
	"testing"

	"onthemat/internal/app/config"
)

func TestPostgres(t *testing.T) {
	c := config.NewConfig()
	c.Load("../../../configs")
	NewPostgresDB(c)
}
