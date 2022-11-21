package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

type AcademyDetailRepsonse struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	CallNumber    string               `json:"callNumber"`
	AddressRoad   string               `json:"addressRoad"`
	AddressDetail *string              `json:"addressDetail"`
	AddressSigun  *string              `json:"addressSigun"`
	CreatedAt     transport.TimeString `json:"createdAt"`
	UpdatedAt     transport.TimeString `json:"updatedAt"`
}

type AcademyListResponse struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	CallNumber    string               `json:"callNumber"`
	AddressRoad   string               `json:"addressRoad"`
	AddressDetail *string              `json:"addressDetail"`
	AddressSigun  *string              `json:"addressSigun"`
	CreatedAt     transport.TimeString `json:"createdAt"`
	UpdatedAt     transport.TimeString `json:"updatedAt"`
}

func NewAcademyListResponse(model []*ent.Academy) []*AcademyListResponse {
	response := make([]*AcademyListResponse, 0)
	for _, v := range model {
		resp := new(AcademyListResponse)
		copier.Copy(&resp, v)
		if v.Edges.Sigungu != nil {
			resp.AddressSigun = &v.Edges.Sigungu.Name
		}
		response = append(response, resp)
	}

	return response
}

func NewAcademyDetailResponse(model *ent.Academy) *AcademyDetailRepsonse {
	resp := new(AcademyDetailRepsonse)
	copier.Copy(&resp, model)
	if model.Edges.Sigungu != nil {
		resp.AddressSigun = &model.Edges.Sigungu.Name
	}
	return resp
}
