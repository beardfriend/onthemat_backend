package request

// ------------------- Create -------------------

// ------------------- Update -------------------

// ------------------- Patch -------------------

type AcademyPatchBody struct {
	Info    *AcademyPatchInfo `json:"info" validate:"dive"`
	YogaIDs *[]int            `json:"yogaIds"`
}

type AcademyPatchInfo struct {
	Name          *string `json:"name"`
	LogoUrl       *string `json:"logoUrl" validate:"omitempty,urlStartHttpHttps"`
	CallNumber    *string `json:"callNumber"`
	AddressRoad   *string `json:"addressRoad"`
	AddressDetail *string `json:"addressDetail"`
	SigunguId     *int    `json:"sigunguId"`
}

// ------------------- List -------------------

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
