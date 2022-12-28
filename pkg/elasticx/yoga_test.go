package elasticx

import (
	"testing"

	"onthemat/internal/app/config"
	"onthemat/internal/app/infrastructure"
)

func TestYogaInit(t *testing.T) {
	c := config.NewConfig()
	c.Load("../../configs")
	ela := infrastructure.NewElasticSearch(c, "../../configs/elastic.crt")
	elax := NewElasticX(ela)

	elax.InitYoga()
}
