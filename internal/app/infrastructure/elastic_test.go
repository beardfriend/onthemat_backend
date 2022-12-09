package infrastructure_test

import (
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
)

func TestElastic(t *testing.T) {
	c := config.NewConfig()
	c.Load("../../../configs")
	infrastructure.NewElasticSearch(c)
}
