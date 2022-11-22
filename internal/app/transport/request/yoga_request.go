package request

// ------------------- Body -------------------
type YogaGroupCreateBody struct {
	Category    string `json:"category" validate:"required"`
	CategoryEng string `json:"categoryEng"`
	Description string `json:"description"`
}

type YogaGroupUpdateBody struct {
	Category    string `json:"category" validate:"required"`
	CategoryEng string `json:"categoryEng"`
	Description string `json:"description"`
}

type YogaCreateBody struct {
	YogaGroupId int     `json:"yogaGroupId"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"string"`
}

type YogaUpdateBody struct {
	YogaGroupId int     `json:"yogaGroupId"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"string"`
}

type YogaRawCreateBody struct {
	YogaName string `json:"name"`
}

// ------------------- Query String -------------------

type YogaGroupListQueries struct {
	PageNo      int     `query:"pageNo"`
	PageSize    int     `query:"pageSize"`
	SearchValue *string `query:"searchValue"`
	SearchKey   *string `query:"searchKey"`
	OrderType   *string `query:"orderType"`
	OrderCol    *string `query:"orderCol"`
}

func NewYogaGroupListQueries() *YogaGroupListQueries {
	return &YogaGroupListQueries{
		PageNo:      1,
		PageSize:    10,
		SearchKey:   nil,
		SearchValue: nil,
	}
}

type YogaGroupDeleteQueries struct {
	Ids []int `query:"ids"`
}

// ------------------- Params -------------------

type YogaUpdateParam struct {
	Id int `json:"id"`
}
