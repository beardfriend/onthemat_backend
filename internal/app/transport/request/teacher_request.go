package request

import "onthemat/internal/app/transport"

// ------------------- Create -------------------
type TeacherCreateBody struct {
	TeacherInfo     TeacherInfoForCreate        `json:"info"`
	YogaIds         []int                       `json:"yogaIds"`
	SigunguIds      []int                       `json:"sigunguIds"`
	WorkExperiences []*WorkExperiencesForCreate `json:"workExperiencs"`
	Certifications  []*CertificationsForCreate  `json:"certifications"`
	YogaRaws        []*YogaRawsForCreate        `json:"yogaRaws"`
}
type WorkExperiencesForCreate struct {
	AcademyName string                `json:"academyName"`
	WorkStartAt transport.TimeString  `json:"workStartAt"`
	WorkEndAt   *transport.TimeString `json:"workEndAt"`
	Description *string               `json:"description"`
}
type CertificationsForCreate struct {
	AgencyName   string                `json:"agencyName"`
	ImageUrl     *string               `json:"imageUrl"`
	ClassStartAt transport.TimeString  `json:"classStartAt"`
	ClassEndAt   *transport.TimeString `json:"classEndAt"`
	Description  *string               `json:"description"`
}

type YogaRawsForCreate struct {
	Name string `json:"name"`
}
type TeacherInfoForCreate struct {
	Name            string  `json:"name"`
	ProfileImageURL *string `json:"profileImageUrl"`
	Age             *int    `json:"age"`
	Introduce       *string `json:"introduce"`
	IsProfileOpen   bool    `json:"isProfileOpen"`
}

// ------------------- Update -------------------

// ___________ body ___________
type TeacherUpdateBody struct {
	TeacherInfo     TeacherInfoForUpdate        `json:"info"`
	YogaIds         []int                       `json:"yogaIds"`
	SigunguIds      []int                       `json:"sigunguIds"`
	WorkExperiences []*WorkExperiencesForUpdate `json:"workExperiencs"`
	Certifications  []*CertificationsForUpdate  `json:"certifications"`
	YogaRaws        []*YogaRawsForPut           `json:"yogaRaws"`
}
type WorkExperiencesForUpdate struct {
	Id          int                   `json:"id"`
	AcademyName string                `json:"academyName"`
	WorkStartAt transport.TimeString  `json:"workStartAt"`
	WorkEndAt   *transport.TimeString `json:"workEndAt"`
	Description *string               `json:"description"`
}
type CertificationsForUpdate struct {
	Id           int                   `json:"id"`
	AgencyName   string                `json:"agencyName"`
	ImageUrl     *string               `json:"imageUrl"`
	ClassStartAt transport.TimeString  `json:"classStartAt"`
	ClassEndAt   *transport.TimeString `json:"classEndAt"`
	Description  *string               `json:"description"`
}

type YogaRawsForPut struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type TeacherInfoForUpdate struct {
	Name            string  `json:"name"`
	ProfileImageURL *string `json:"profileImageUrl"`
	Age             *int    `json:"age"`
	Introduce       *string `json:"introduce"`
	IsProfileOpen   bool    `json:"isProfileOpen"`
}

// ___________ param ___________
type TeacherUpdateParam struct {
	Id int `params:"id" validate:"required"`
}

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
