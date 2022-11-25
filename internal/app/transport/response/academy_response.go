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
	AddressSigun  string               `json:"addressSigun"`
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
	AddressSigun  string               `json:"addressSigun"`
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

		resp.AddressSigun = v.Edges.AreaSigungu.Name

		if len(v.Edges.Yoga) > 0 {
			copier.Copy(&resp.Yoga, v.Edges.Yoga)
		} else {
			resp.Yoga = make([]yoga, 0)
		}
		response = append(response, resp)
	}

	return response
}

func NewAcademyDetailResponse(m *ent.Academy) *AcademyDetailRepsonse {
	resp := &AcademyDetailRepsonse{
		ID:            m.ID,
		Name:          m.Name,
		CallNumber:    m.CallNumber,
		AddressRoad:   m.AddressRoad,
		AddressDetail: m.AddressDetail,
		AddressSigun:  m.Edges.AreaSigungu.Name,
	}

	if len(m.Edges.Yoga) > 0 {
		for _, v := range m.Edges.Yoga {
			resp.Yoga = append(resp.Yoga, yoga{
				ID:      v.ID,
				NameKor: v.NameKor,
			})
		}
	} else {
		resp.Yoga = make([]yoga, 0)
	}

	return resp
}
