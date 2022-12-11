package elastic

import "github.com/elastic/go-elasticsearch/v8"

type ElasticX struct {
	ela *elasticsearch.Client
}

func NewElasticX(ela *elasticsearch.Client) *ElasticX {
	return &ElasticX{
		ela: ela,
	}
}

func (e *ElasticX) AuthMigration() {
	e.InitYoga()
}
