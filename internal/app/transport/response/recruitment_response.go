package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
)

type RecruitmentListReponse struct {
	ID             int                    `json:"id"`
	AcademyName    string                 `json:"academyName"`
	Yogas          []string               `json:"yogas"`
	Sigungu        string                 `json:"sigungu"`
	StartDateTimes []transport.TimeString `json:"startDateTimes"`
}

func NewRecruitmentListResponse(model []*ent.Recruitment) []*RecruitmentListReponse {
	response := make([]*RecruitmentListReponse, 0)

	for _, v := range model {
		yogas := make([]string, 0)
		startDateTimes := make([]transport.TimeString, 0)
		for _, j := range v.Edges.RecruitmentInstead {
			for _, y := range j.Edges.Yoga {
				yogas = append(yogas, y.NameKor)
			}
			for _, s := range j.Schedule {
				startDateTimes = append(startDateTimes, s.StartDateTime)
			}

		}

		response = append(response, &RecruitmentListReponse{
			ID:             v.ID,
			AcademyName:    v.Edges.Writer.Name,
			Yogas:          yogas,
			Sigungu:        v.Edges.Writer.Edges.AreaSigungu.Name,
			StartDateTimes: startDateTimes,
		})
	}

	return response
}
