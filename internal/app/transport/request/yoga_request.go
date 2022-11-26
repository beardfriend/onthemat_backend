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
	Description *string `json:"description"`
}

type YogaUpdateBody struct {
	YogaGroupId int     `json:"yogaGroupId" validate:"required"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"description"`
}

type YogaPatcheBody struct {
	YogaGroupId *int    `json:"yogaGroupId,omitempty" validate:"required"`
	NameKor     *string `json:"nameKor,omitempty"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level,omitempty"`
	Description *string `json:"description,omitempty"`
}

type YogaGroupsDeleteBody struct {
	Ids []int `json:"ids" validate:"required"`
}

type YogaRawCreateBody struct {
	YogaName string `json:"name" validate:"required"`
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

type YogaListQuery struct {
	GroupId int `query:"groupId"`
}

type YogaRawListQuery struct {
	UserId int `query:"userId"`
}

// ------------------- Params -------------------

type YogaUpdateParam struct {
	Id int `params:"id"`
}

type YogaDeleteParam struct {
	Id int `params:"id"`
}

type YogaDeleteRawParam struct {
	Id int `params:"id"`
}
