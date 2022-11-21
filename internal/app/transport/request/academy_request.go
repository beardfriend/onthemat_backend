package request

type AcademyListQueries struct {
	PageNo      int     `query:"pageNo"`
	PageSize    int     `query:"pageSize"`
	SearchValue *string `query:"searchValue"`
	SearchKey   *string `query:"searchKey"`
	OrderType   *string `query:"orderType"`
	OrderCol    *string `query:"orderCol"`
}

func NewAcademyListQueries() *AcademyListQueries {
	return &AcademyListQueries{
		PageNo:      1,
		PageSize:    10,
		SearchKey:   nil,
		SearchValue: nil,
		OrderType:   nil,
		OrderCol:    nil,
	}
}
