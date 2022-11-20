package transport

import (
	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

// ------------------- Request -------------------

// ___________ Param ___________

type AcademyDetailParam struct {
	Id int `params:"id" validate:"required,numeric"`
}

type AcademyUpdateParam struct {
	Id int `params:"id" validate:"required,numeric"`
}

// ___________ Query ___________

type AcademyListQueries struct {
	PageNo      int     `query:"pageNo"`
	PageSize    int     `query:"pageSize"`
	SearchKey   *string `query:"searchKey"`
	SearchValue *string `query:"searchValue"`
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

// ___________ Body ___________

type AcademyCreateRequestBody struct {
	Name          string `json:"name" validate:"required"`
	LogoUrl       string `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	BusinessCode  string `json:"businessCode" validate:"required"`
	CallNumber    string `json:"callNumber" validate:"required,phoneNumNoDash"`
	AddressRoad   string `json:"addressRoad" validate:"required"`
	AddressSigun  string `json:"addressSiGun" validate:"required"`
	AddressGu     string `json:"addressGu" validate:"required"`
	AddressDong   string `json:"addressDong" validate:"required"`
	AddressDetail string `json:"addressDetail"`
	AddressX      string `json:"addressX" validate:"required"`
	AddressY      string `json:"addressY" validate:"required"`
}

type AcademyUpdateRequestBody struct {
	Name          string `json:"name" validate:"required"`
	LogoUrl       string `json:"logoUrl" validate:"required,urlStartHttpHttps"`
	CallNumber    string `json:"callNumber"`
	AddressRoad   string `json:"addressRoad" validate:"required"`
	AddressDetail string `json:"addressDetail" validate:"required"`
	AddressSigun  string `json:"addressSigun" validate:"required"`
	AddressGu     string `json:"addressGu" validate:"required"`
	AddressDong   string `json:"addressDong" validate:"required"`
	AddressX      string `json:"addressX" validate:"required"`
	AddressY      string `json:"addressY" validate:"required"`
}

// ------------------- Response -------------------

type AcademyDetailRepsonse struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	CallNumber    string     `json:"callNumber"`
	AddressRoad   string     `json:"addressRoad"`
	AddressDetail string     `json:"addressDetail"`
	AddressSigun  string     `json:"addressSigun"`
	AddressGu     string     `json:"addressGu"`
	AddressX      string     `json:"addressX"`
	AddressY      string     `json:"addressY"`
	CreatedAt     TimeString `json:"createdAt"`
	UpdatedAt     TimeString `json:"updatedAt"`
}

type AcademyListResponse struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	CallNumber    string     `json:"callNumber"`
	AddressRoad   string     `json:"addressRoad"`
	AddressDetail string     `json:"addressDetail"`
	AddressSigun  string     `json:"addressSigun"`
	AddressGu     string     `json:"addressGu"`
	CreatedAt     TimeString `json:"createdAt"`
	UpdatedAt     TimeString `json:"updatedAt"`
}

func NewAcademyListResponse(model []*ent.Academy) []*AcademyListResponse {
	response := make([]*AcademyListResponse, 0)
	for _, v := range model {
		resp := new(AcademyListResponse)
		copier.Copy(&resp, v)
		response = append(response, resp)
	}

	return response
}

func NewAcademyDetailResponse(model *ent.Academy) *AcademyDetailRepsonse {
	resp := new(AcademyDetailRepsonse)
	copier.Copy(&resp, model)
	return resp
}
