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

type ElasticDeleteUpdateQueryResponse struct {
	Took                 int           `json:"took"`
	TimedOut             bool          `json:"timed_out"`
	Total                int           `json:"total"`
	Updated              int           `json:"updated"`
	Deleted              int           `json:"deleted"`
	Batches              int           `json:"batches"`
	VersionConflicts     int           `json:"version_conflicts"`
	Noops                int           `json:"noops"`
	Retries              Retries       `json:"retries"`
	ThrottledMillis      int           `json:"throttled_millis"`
	RequestsPerSecond    interface{}   `json:"requests_per_second"`
	ThrottledUntilMillis int           `json:"throttled_until_millis"`
	Failures             []interface{} `json:"failures"`
}

type Retries struct {
	Bulk   int `json:"bulk"`
	Search int `json:"search"`
}
