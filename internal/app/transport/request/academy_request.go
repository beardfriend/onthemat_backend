package request

type AcademyListQueries struct {
	PageNo      int     `query:"pageNo"`
	PageSize    int     `query:"pageSize"`
	AcademyName *string `query:"academyName"`
	YogaIDs     *[]int  `query:"yogaIDs"`
	SigunGuID   *int    `query:"sigunguID"`
	OrderType   *string `query:"orderType"`
	OrderCol    *string `query:"orderCol"`
}

func NewAcademyListQueries() *AcademyListQueries {
	return &AcademyListQueries{
		PageNo:      1,
		PageSize:    10,
		AcademyName: nil,
		YogaIDs:     nil,
		SigunGuID:   nil,
		OrderType:   nil,
		OrderCol:    nil,
	}
}

type AddYogaBody struct {
	Ids []int `json:"yogaIDs"`
}

type DeleteYogaBody struct {
	Ids []int `json:"yogaIDs"`
}
