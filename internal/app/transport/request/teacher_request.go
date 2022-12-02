package request

import "onthemat/internal/app/transport"

// ------------------- Create -------------------
type (
	TeacherCreateBody struct {
		TeacherInfo     TeacherInfoForCreate        `json:"info" validate:"required,dive"`
		YogaIds         []int                       `json:"yogaIds"`
		SigunguIds      []int                       `json:"sigunguIds"`
		WorkExperiences []*WorkExperiencesForCreate `json:"workExperiences" validate:"omitempty,dive"`
		Certifications  []*CertificationsForCreate  `json:"certifications" validate:"omitempty,dive"`
		YogaRaws        []*YogaRawsForCreate        `json:"yogaRaws" validate:"omitempty,dive"`
	}
	WorkExperiencesForCreate struct {
		AcademyName string                `json:"academyName" validate:"omitempty,required"`
		WorkStartAt transport.TimeString  `json:"workStartAt" validate:"omitempty,required"`
		WorkEndAt   *transport.TimeString `json:"workEndAt"`
		Description *string               `json:"description"`
	}
	CertificationsForCreate struct {
		AgencyName   string                `json:"agencyName" validate:"omitempty,required"`
		ImageUrl     *string               `json:"imageUrl"`
		ClassStartAt transport.TimeString  `json:"classStartAt" validate:"omitempty,required"`
		ClassEndAt   *transport.TimeString `json:"classEndAt"`
		Description  *string               `json:"description"`
	}

	YogaRawsForCreate struct {
		Name string `json:"name"`
	}
	TeacherInfoForCreate struct {
		Name            string  `json:"name" validate:"omitempty,required"`
		ProfileImageURL *string `json:"profileImageUrl"`
		Age             *int    `json:"age"`
		Introduce       *string `json:"introduce"`
		IsProfileOpen   bool    `json:"isProfileOpen"`
	}
)

// ------------------- Update -------------------

// ___________ body ___________
type (
	TeacherUpdateBody struct {
		TeacherInfo     TeacherInfoForUpdate        `json:"info" validate:"omitempty,dive"`
		YogaIds         []int                       `json:"yogaIds"`
		SigunguIds      []int                       `json:"sigunguIds"`
		WorkExperiences []*WorkExperiencesForUpdate `json:"workExperiences" validate:"omitempty,dive"`
		Certifications  []*CertificationsForUpdate  `json:"certifications" validate:"omitempty,dive"`
		YogaRaws        []*YogaRawsForUpdate        `json:"yogaRaws" validate:"omitempty,dive"`
	}
	WorkExperiencesForUpdate struct {
		Id          int                   `json:"id" validate:"omitempty,required"`
		AcademyName string                `json:"academyName" validate:"omitempty,required"`
		WorkStartAt transport.TimeString  `json:"workStartAt" validate:"omitempty,required"`
		WorkEndAt   *transport.TimeString `json:"workEndAt"`
		Description *string               `json:"description"`
	}
	CertificationsForUpdate struct {
		Id           int                   `json:"id" validate:"omitempty,required"`
		AgencyName   string                `json:"agencyName" validate:"omitempty,required"`
		ImageUrl     *string               `json:"imageUrl"`
		ClassStartAt transport.TimeString  `json:"classStartAt" validate:"omitempty,required"`
		ClassEndAt   *transport.TimeString `json:"classEndAt"`
		Description  *string               `json:"description"`
	}

	YogaRawsForUpdate struct {
		Id   int    `json:"id" validate:"omitempty,required"`
		Name string `json:"name" validate:"omitempty,required"`
	}
	TeacherInfoForUpdate struct {
		Name            string  `json:"name" validate:"required"`
		ProfileImageURL *string `json:"profileImageUrl"`
		Age             *int    `json:"age"`
		Introduce       *string `json:"introduce"`
		IsProfileOpen   bool    `json:"isProfileOpen"`
	}
)

// ___________ param ___________
type TeacherUpdateParam struct {
	Id int `params:"id" validate:"required"`
}

// ------------------- Patch -------------------

type (
	TeacherPatchBody struct {
		YogaIds         *[]int                     `json:"yogaIds"`
		SigunguIds      *[]int                     `json:"sigunguIds"`
		WorkExperiences []*WorkExperiencesForPatch `json:"workExperiencs" validate:"omitempty,dive"`
		Certifications  []*CertificationsForPatch  `json:"certifications" validate:"omitempty,dive"`
		YogaRaws        []*YogaRawsForPatch        `json:"yogaRaws" validate:"omitempty,dive"`
		TeacherInfo     *TeacherInfoForPatch       `json:"info" validate:"omitempty,dive"`
	}

	WorkExperiencesForPatch struct {
		Id          *int                  `json:"id"`
		AcademyName *string               `json:"academyName" validate:"must"`
		WorkStartAt *transport.TimeString `json:"workStartAt" validate:"must"`
		WorkEndAt   *transport.TimeString `json:"workEndAt"`
		Description *string               `json:"description"`
	}

	CertificationsForPatch struct {
		Id           *int                  `json:"id"`
		AgencyName   *string               `json:"agencyName" validate:"must"`
		ImageUrl     *string               `json:"imageUrl"`
		ClassStartAt *transport.TimeString `json:"classStartAt" validate:"must"`
		ClassEndAt   *transport.TimeString `json:"classEndAt"`
		Description  *string               `json:"description"`
	}

	YogaRawsForPatch struct {
		Id   *int    `json:"id"`
		Name *string `json:"name" validate:"must"`
	}
	TeacherInfoForPatch struct {
		ProfileImageURL *string `json:"profileImageUrl"`
		Name            *string `json:"name"`
		Age             *int    `json:"age"`
		IsProfileOpen   *bool   `json:"isProfileOpen"`
		Introduce       *string `json:"introduce"`
	}
)

type TeacherPatchParam struct {
	Id int `params:"id" validate:"required"`
}
