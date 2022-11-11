package common

type ListParams struct {
	PageNo      int
	PageSize    int
	SearchKey   *string
	SearchValue *string
	OrderType   *string
	OrderCol    *string
}

type TotalParams struct {
	SearchKey   *string
	SearchValue *string
}
