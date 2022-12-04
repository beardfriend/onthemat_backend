package response

import (
	"onthemat/internal/app/model"
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
)

// ------------------- List -------------------

type RecruitmentListReponse struct {
	ID             int                    `json:"id"`
	AcademyName    string                 `json:"academyName"`
	Yogas          []string               `json:"yogas"`
	Sigungu        string                 `json:"sigungu"`
	StartDateTimes []transport.TimeString `json:"startDateTimes"`
	CreatedAt      transport.TimeString   `json:"createdAt"`
	UpdatedAt      transport.TimeString   `json:"updatedAt"`
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
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	return response
}

// ------------------- Get -------------------

type RecruitmentReponse struct {
	ID          int                  `json:"id"`
	IsFinish    bool                 `json:"isFinish"`
	InsteadInfo []*insteadInfo       `json:"insteadInfo"`
	AcademyInfo *academyInfo         `json:"academyInfo"`
	CreatedAt   transport.TimeString `json:"createdAt"`
	UpdatedAt   transport.TimeString `json:"updatedAt"`
}

type insteadInfo struct {
	ID             int                  `json:"id"`
	MinCareer      string               `json:"minCareer"`
	Pay            string               `json:"pay"`
	PasserId       *int                 `json:"passerId"`
	ApplicantCount int                  `json:"applicantCount"`
	Schedules      []*model.Schedule    `json:"schedules"`
	Yogas          []*yogaInfo          `json:"yogas"`
	CreatedAt      transport.TimeString `json:"createdAt"`
	UpdatedAt      transport.TimeString `json:"updatedAt"`
}

type academyInfo struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	AddressRoad   string  `json:"addressRoad"`
	AddressDetail *string `json:"addressDetail"`
}

type yogaInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewRecruitmentResponse(v *ent.Recruitment) *RecruitmentReponse {
	riInfo := make([]*insteadInfo, 0)
	for _, v := range v.Edges.RecruitmentInstead {
		yogas := make([]*yogaInfo, 0)

		for _, y := range v.Edges.Yoga {
			yogas = append(yogas, &yogaInfo{
				Id:   y.ID,
				Name: y.NameKor,
			})
		}

		riInfo = append(riInfo, &insteadInfo{
			ID:             v.ID,
			MinCareer:      v.MinCareer,
			PasserId:       v.TeacherID,
			ApplicantCount: len(v.Edges.Applicant),
			Pay:            v.Pay,
			Schedules:      v.Schedule,
			Yogas:          yogas,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		})
	}

	academy := v.Edges.Writer
	resp := &RecruitmentReponse{
		ID:          v.ID,
		IsFinish:    v.IsFinish,
		InsteadInfo: riInfo,
		AcademyInfo: &academyInfo{
			Id:            academy.ID,
			Name:          academy.Name,
			AddressRoad:   academy.AddressRoad,
			AddressDetail: academy.AddressDetail,
		},
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}

	return resp
}
