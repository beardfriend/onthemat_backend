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
	ID            int    `json:"id"`
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
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CallNumber    string `json:"call_number"`
	AddressRoad   string `json:"address_road"`
	AddressDetail string `json:"address_detail"`
	AddressSigun  string `json:"address_sigun"`
	AddressGu     string `json:"address_gu"`
	AddressX      string `json:"address_x"`
	AddressY      string `json:"address_y"`
}

func NewAcademyDetailResponse(model *ent.Acadmey) *AcademyDetailRepsonse {
	resp := new(AcademyDetailRepsonse)
	copier.Copy(&resp, model)
	return resp
}
