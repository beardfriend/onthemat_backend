package request

import "onthemat/internal/app/transport"

// ------------------- Patch -------------------

type TeacherPatchBody struct {
	YogaIds         *[]int                     `json:"yogaIds"`
	SigunguIds      *[]int                     `json:"sigunguIds"`
	WorkExperiences []*WorkExperiencesForPatch `json:"workExperiencs"`
	Certifications  []*CertificationsForPatch  `json:"certifications"`
	YogaRaws        []*YogaRawsForPatch        `json:"yogaRaws"`
	TeacherInfo     *TeacherInfoForPatch       `json:"teacherInfo"`
}

type WorkExperiencesForPatch struct {
	Id          *int                  `json:"id"`
	AcademyName *string               `json:"academyName"`
	WorkStartAt *transport.TimeString `json:"workStartAt"`
	WorkEndAt   *transport.TimeString `json:"workEndAt"`
	Description *string               `json:"description"`
}

type CertificationsForPatch struct {
	Id           *int                  `json:"id"`
	AgencyName   *string               `json:"agencyName"`
	ImageUrl     *string               `json:"imageUrl"`
	ClassStartAt *transport.TimeString `json:"classStartAt"`
	ClassEndAt   *transport.TimeString `json:"classEndAt"`
	Description  *string               `json:"description"`
}

type YogaRawsForPatch struct {
	Name *string `json:"name"`
}
type TeacherInfoForPatch struct {
	ProfileImageURL *string `json:"profileImageUrl"`
	Name            *string `json:"name"`
	Age             *int    `json:"age"`
	IsProfileOpen   *bool   `json:"isProfileOpen"`
	Introduce       *string `json:"introduce"`
}
