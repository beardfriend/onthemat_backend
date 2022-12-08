package infrastructure_test

import (
	"testing"

	"onthemat/internal/app/infrastructure"
)

func TestElastic(t *testing.T) {
	infrastructure.NewElasticSearch()
}
