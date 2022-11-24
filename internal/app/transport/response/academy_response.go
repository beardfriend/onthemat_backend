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
	Yoga          []yoga               `json:"yoga"`
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
	Yoga          []yoga               `json:"yoga"`
	CreatedAt     transport.TimeString `json:"createdAt"`
	UpdatedAt     transport.TimeString `json:"updatedAt"`
}

type yoga struct {
	ID      int    `json:"id"`
	NameKor string `json:"nameKor"`
}

func NewAcademyListResponse(model []*ent.Academy) []*AcademyListResponse {
	response := make([]*AcademyListResponse, 0)
	for _, v := range model {
		resp := new(AcademyListResponse)
		copier.Copy(&resp, v)
		if v.Edges.AreaSigungu != nil {
			resp.AddressSigun = &v.Edges.AreaSigungu.Name
		}

		if len(v.Edges.Yoga) > 0 {
			copier.Copy(&resp.Yoga, v.Edges.Yoga)
		} else {
			resp.Yoga = make([]yoga, 0)
		}
		response = append(response, resp)
	}

	return response
}

func NewAcademyDetailResponse(model *ent.Academy) *AcademyDetailRepsonse {
	resp := new(AcademyDetailRepsonse)
	copier.Copy(&resp, model)

	if model.Edges.AreaSigungu != nil {
		resp.AddressSigun = &model.Edges.AreaSigungu.Name
	}

	if len(model.Edges.Yoga) > 0 {
		copier.Copy(&resp.Yoga, model.Edges.Yoga)
	} else {
		resp.Yoga = make([]yoga, 0)
	}

	return resp
}
