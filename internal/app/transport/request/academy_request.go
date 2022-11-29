package request

// ------------------- Create -------------------

type AcademyCreateBody struct {
	Info    AcademyInfoForCreate `json:"info" validate:"required,dive,required"`
	YogaIDs []int                `json:"yogaIDs" validate:"dive,required"`
}

type AcademyInfoForCreate struct {
	SigunguID     int     `json:"sigunguId" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	LogoUrl       string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber    string  `json:"callNumber"`
	BusinessCode  string  `json:"businessCode" validate:"required"`
	AddressRoad   string  `json:"addressRoad" validate:"required"`
	AddressDetail *string `json:"addressDetail"`
}

// ------------------- Update -------------------

type AcademyUpdateBody struct {
	Info    AcademyInfoForUpdate `json:"info" validate:"required,dive,required"`
	YogaIDs []int                `json:"yogaIds" validate:"dive,required"`
}

type AcademyInfoForUpdate struct {
	SigunguID     int     `json:"sigunguId" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	LogoUrl       string  `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber    string  `json:"callNumber"`
	AddressRoad   string  `json:"addressRoad" validate:"required"`
	AddressDetail *string `json:"addressDetail"`
}

type AcademyUpdateParam struct {
	Id int `params:"id" validate:"required"`
}

// ------------------- Patch -------------------

type AcademyPatchBody struct {
	Info    *AcademyInfoForPatch `json:"info" validate:"dive"`
	YogaIDs *[]int               `json:"yogaIds"`
}

type AcademyInfoForPatch struct {
	Name          *string `json:"name"`
	LogoUrl       *string `json:"logoUrl" validate:"omitempty,urlStartHttpHttps"`
	CallNumber    *string `json:"callNumber"`
	AddressRoad   *string `json:"addressRoad"`
	AddressDetail *string `json:"addressDetail"`
	SigunguId     *int    `json:"sigunguId"`
}
type AcademyPatchParam struct {
	Id int `params:"id" validate:"required"`
}

// ------------------- Get -------------------
type AcademyDetailParam struct {
	Id int `params:"id" validate:"required,numeric"`
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
