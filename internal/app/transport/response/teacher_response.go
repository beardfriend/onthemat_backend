package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"
)

type TeacherResponse struct {
	Id                  int                   `json:"id"`
	Name                string                `json:"name"`
	Age                 *int                  `json:"age"`
	Introduce           *string               `json:"introduce"`
	ProfileImageUrl     *string               `json:"profileImageUrl"`
	IsProfileOpen       bool                  `json:"isProfileOpen"`
	Yoga                []yoga                `json:"yoga"`
	PossibleWorkSigungu []possibleWorkSigungu `json:"possibleWorkSigungu"`
	Certifications      []certifications      `json:"certifications"`
	WorkExperiences     []workExperiences     `json:"workExperiences"`
	CreatedAt           transport.TimeString  `json:"createdAt"`
	UpdatedAt           transport.TimeString  `json:"updatedAt"`
}

type possibleWorkSigungu struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type certifications struct {
	Id           int                   `json:"id"`
	AgencyName   string                `json:"agencyName"`
	ImageUrl     *string               `json:"imageUrl"`
	ClassStartAt transport.TimeString  `json:"classStartAt"`
	ClassEndAt   *transport.TimeString `json:"classEndAt"`
	Description  *string               `json:"description"`
	CreatedAt    transport.TimeString  `json:"createdAt"`
	UpdatedAt    transport.TimeString  `json:"updatedAt"`
}

type workExperiences struct {
	Id          int                   `json:"id"`
	AcademyName string                `json:"academyName"`
	WorkStartAt transport.TimeString  `json:"workStartAt"`
	WorkEndAt   *transport.TimeString `json:"workEndAt"`
	Description *string               `json:"description"`
	CreatedAt   transport.TimeString  `json:"createdAt"`
	UpdatedAt   transport.TimeString  `json:"updatedAt"`
}

func NewTeacherResponse(d *ent.Teacher) *TeacherResponse {
	// info
	resp := &TeacherResponse{
		Id:              d.ID,
		Name:            d.Name,
		Age:             d.Age,
		Introduce:       d.Introduce,
		ProfileImageUrl: d.ProfileImageUrl,
		IsProfileOpen:   d.IsProfileOpen,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
	resp.Yoga = make([]yoga, 0)

	// yoga
	index := 0
	for _, v := range d.Edges.Yoga {
		resp.Yoga = append(resp.Yoga, yoga{
			Index:       index,
			ID:          v.ID,
			NameKor:     v.NameKor,
			IsReference: true,
		})
		index++
	}
	for _, v := range d.Edges.YogaRaw {
		resp.Yoga = append(resp.Yoga, yoga{
			Index:       index,
			ID:          v.ID,
			NameKor:     v.Name,
			IsReference: false,
		})
		index++
	}

	// sigungu
	resp.PossibleWorkSigungu = make([]possibleWorkSigungu, 0)
	for _, v := range d.Edges.Sigungu {
		resp.PossibleWorkSigungu = append(resp.PossibleWorkSigungu, possibleWorkSigungu{
			Id:   v.ID,
			Name: v.Name,
		})
	}

	// experience
	resp.WorkExperiences = make([]workExperiences, 0)
	for _, v := range d.Edges.WorkExperience {
		resp.WorkExperiences = append(resp.WorkExperiences, workExperiences{
			Id:          v.ID,
			AcademyName: v.AcademyName,
			WorkStartAt: v.WorkStartAt,
			WorkEndAt:   v.WorkEndAt,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	// certification
	resp.Certifications = make([]certifications, 0)
	for _, v := range d.Edges.Certification {
		resp.Certifications = append(resp.Certifications, certifications{
			Id:           v.ID,
			AgencyName:   v.AgencyName,
			ImageUrl:     v.ImageUrl,
			ClassStartAt: v.ClassStartAt,
			ClassEndAt:   v.ClassEndAt,
			Description:  v.Description,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	return resp
}
