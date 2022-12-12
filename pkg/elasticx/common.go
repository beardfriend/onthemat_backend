package elasticx

type ElasticSearchResponse[T any] struct {
	Hits struct {
		Hits []T `json:"hits"`
	} `json:"hits"`
}

type ElasticSearchListBody[T any] struct {
	Index  string  `json:"_index"`
	Id     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
}
