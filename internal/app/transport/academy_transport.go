package transport

import (
	"time"

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
	SearchKey   *string `query:"searchKey" validate:"omitempty,oneof=name gu NAME GU"`
	SearchValue *string `query:"searchValue"`
	OrderType   *string `query:"orderType"`
	OrderCol    *string `query:"orderCol" `
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
	Name          string `json:"name"`
	LogoUrl       string `json:"logoUrl"`
	BusinessCode  string `json:"businessCode"`
	CallNumber    string `json:"callNumber"`
	AddressRoad   string `json:"addressRoad"`
	AddressSigun  string `json:"addressSiGun"`
	AddressGu     string `json:"addressGu"`
	AddressDong   string `json:"addressDong"`
	AddressDetail string `json:"addressDetail"`
	AddressX      string `json:"addressX"`
	AddressY      string `json:"addressY"`
}

type AcademyUpdateRequestBody struct {
	Name          string `json:"name"`
	CallNumber    string `json:"callNumber"`
	AddressRoad   string `json:"addressRoad"`
	AddressDetail string `json:"addressDetail"`
	AddressSigun  string `json:"addressSigun"`
	AddressGu     string `json:"addressGu"`
	AddressDong   string `json:"addressDong"`
	AddressX      string `json:"addressX"`
	AddressY      string `json:"addressY"`
}

// ------------------- Response -------------------

type AcademyDetailRepsonse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	CallNumber    string    `json:"callNumber"`
	AddressRoad   string    `json:"addressRoad"`
	AddressDetail string    `json:"addressDetail"`
	AddressSigun  string    `json:"addressSigun"`
	AddressGu     string    `json:"addressGu"`
	AddressX      string    `json:"addressX"`
	AddressY      string    `json:"addressY"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type AcademyListResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	CallNumber    string    `json:"callNumber"`
	AddressRoad   string    `json:"addressRoad"`
	AddressDetail string    `json:"addressDetail"`
	AddressSigun  string    `json:"addressSigun"`
	AddressGu     string    `json:"addressGu"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
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
