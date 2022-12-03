package request

import (
	"onthemat/internal/app/transport"
)

// ------------------- Create -------------------

// ___________ body ___________

type (
	RecruitmentCreateBody struct {
		Info        RecruitmentInfoForCreate       `json:"info" validate:"required,dive"`
		InsteadInfo []*RecruitmentInsteadForCreate `json:"insteadInfo" validate:"omitempty,dive"`
	}

	RecruitmentInfoForCreate struct {
		IsOpen bool `json:"isOpen"`
	}
	RecruitmentInsteadForCreate struct {
		MinCareer string     `json:"minCareer" validate:"required"`
		Pay       string     `json:"pay" validate:"required"`
		Schedules []Schedule `json:"schedules" validate:"required,dive"`
	}

	Schedule struct {
		StartDateTime transport.TimeString `json:"startDateTime" valiate:"required"`
		EndDateTime   transport.TimeString `json:"endDateTime" valiate:"required"`
	}
)

// ------------------- Update -------------------

// ___________ body ___________

type (
	RecruitmentUpdateBody struct {
		Info        RecruitmentInfoForUpdate       `json:"info" validate:"required,dive"`
		InsteadInfo []*RecruitmentInsteadForUpdate `json:"insteadInfo" validate:"omitempty,dive"`
	}

	RecruitmentInfoForUpdate struct {
		IsFinish bool `json:"isFinish"`
		IsOpen   bool `json:"isOpen"`
	}
	RecruitmentInsteadForUpdate struct {
		Id        int        `json:"id"`
		MinCareer string     `json:"minCareer" validate:"required"`
		Pay       string     `json:"pay" validate:"required"`
		Schedules []Schedule `json:"schedules" validate:"required,dive"`
	}
)

// ___________ Param ___________

type RecruitmentUpdateParam struct {
	Id int `param:"id" validate:"required"`
}

// ------------------- Patch -------------------

// ___________ body ___________

type (
	RecruitmentPatchBody struct {
		Info        *RecruitmentInfoForPatch      `json:"info" validate:"omitempty,dive"`
		InsteadInfo []*RecruitmentInsteadForPatch `json:"instead_info" validate:"omitempty,dive"`
	}

	RecruitmentInfoForPatch struct {
		ID       *int  `json:"id"`
		IsOpen   *bool `json:"isOpen"`
		IsFinish *bool `json:"isFinish"`
	}
	RecruitmentInsteadForPatch struct {
		ID            *int                  `json:"id"`
		MinCareer     *string               `json:"minCareer" validate:"must"`
		Pay           *string               `json:"pay" validate:"must"`
		StartDateTime *transport.TimeString `json:"startDateTime" validate:"must"`
		EndDateTime   *transport.TimeString `json:"endDateTime" validate:"must"`
		PasserId      *int                  `json:"passerId"`
	}
)

// ___________ Param ___________

type RecruitmentPatchParam struct {
	Id int `param:"id" validate:"required"`
}
