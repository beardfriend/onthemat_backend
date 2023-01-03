package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

type AcademyRepsonse struct {
	ID             int                  `json:"id"`
	Name           string               `json:"name"`
	LogoUrl        string               `json:"logoUrl"`
	CallNumber     string               `json:"callNumber"`
	AddressRoad    string               `json:"addressRoad"`
	AddressDetail  *string              `json:"addressDetail"`
	AddressSigungu string               `json:"addressSigun"`
	SigunguId      int                  `json:"sigunguId"`
	Yoga           []yoga               `json:"yoga"`
	CreatedAt      transport.TimeString `json:"createdAt"`
	UpdatedAt      transport.TimeString `json:"updatedAt"`
}

type yoga struct {
	Index       int    `json:"index"`
	ID          int    `json:"id"`
	NameKor     string `json:"nameKor"`
	IsReference bool   `json:"isReference"`
}

func NewAcademyListResponse(model []*ent.Academy) []*AcademyRepsonse {
	response := make([]*AcademyRepsonse, 0)
	for _, v := range model {
		resp := new(AcademyRepsonse)
		copier.Copy(&resp, v)

		resp.AddressSigungu = v.Edges.AreaSigungu.Name

		if len(v.Edges.Yoga) > 0 {
			copier.Copy(&resp.Yoga, v.Edges.Yoga)
		} else {
			resp.Yoga = make([]yoga, 0)
		}
		response = append(response, resp)
	}

	return response
}

func NewAcademyDetailResponse(m *ent.Academy) *AcademyRepsonse {
	resp := &AcademyRepsonse{
		ID:             m.ID,
		Name:           m.Name,
		LogoUrl:        m.LogoUrl,
		CallNumber:     m.CallNumber,
		AddressRoad:    m.AddressRoad,
		AddressDetail:  m.AddressDetail,
		AddressSigungu: m.Edges.AreaSigungu.Name,
		SigunguId:      m.SigunguID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if len(m.Edges.Yoga) > 0 {
		for i, v := range m.Edges.Yoga {
			resp.Yoga = append(resp.Yoga, yoga{
				Index:       i,
				ID:          v.ID,
				NameKor:     v.NameKor,
				IsReference: true,
			})
		}
	} else {
		resp.Yoga = make([]yoga, 0)
	}

	return resp
}
