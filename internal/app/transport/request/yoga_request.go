package request

import (
	"onthemat/internal/app/utils"
)

// ------------------- Yoga -------------------

// ___________ Create ___________
type YogaCreateBody struct {
	YogaGroupId int     `json:"yogaGroupId" validate:"required"`
	NameKor     string  `json:"nameKor" validate:"required"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"description"`
}

// ___________ Update ___________
type YogaUpdateBody struct {
	YogaGroupId int     `json:"yogaGroupId" validate:"required"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"description"`
}

type YogaUpdateParam struct {
	Id int `params:"id"`
}

// ___________ Patch ___________
type YogaPatchBody struct {
	YogaGroupId *int    `json:"yogaGroupId,omitempty" validate:"required"`
	NameKor     *string `json:"nameKor,omitempty"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level,omitempty"`
	Description *string `json:"description,omitempty"`
}

type YogaPatchParam struct {
	Id int `params:"id"`
}

// ___________ Delete ___________

type YogaDeleteParam struct {
	Id int `params:"id"`
}

// ___________ List ___________

type YogaListQuery struct {
	GroupId int `query:"groupId"`
}

// ------------------- YogaGroup -------------------

// ___________ Create ___________
type YogaGroupCreateBody struct {
	Category    string  `json:"category" validate:"required"`
	CategoryEng string  `json:"categoryEng"`
	Description *string `json:"description"`
}

// ___________ Update ___________

type YogaGroupUpdateBody struct {
	Category    string  `json:"category" validate:"required"`
	CategoryEng string  `json:"categoryEng"`
	Description *string `json:"description"`
}

// ___________ Patch ___________
type YogaGroupPatchBody struct {
	Category    *string `json:"category"`
	CategoryEng *string `json:"categoryEng"`
	Description *string `json:"description"`
}

type YogaGroupPatchParam struct {
	Id int `params:"id"`
}

// ___________ Delete ___________

type YogaGroupDeleteParam struct {
	Ids []int `params:"ids"`
}

// ___________ List ___________

type YogaGroupListQueries struct {
	PageNo    int     `query:"pageNo"`
	PageSize  int     `query:"pageSize"`
	Category  *string `query:"category"`
	OrderType *string `query:"orderType"`
}

func NewYogaGroupListQueries() *YogaGroupListQueries {
	return &YogaGroupListQueries{
		PageNo:    1,
		PageSize:  10,
		Category:  nil,
		OrderType: utils.String("DESC"),
	}
}

// ------------------- Yoga Raw -------------------

type YogaRawCreateBody struct {
	YogaName string `json:"name" validate:"required"`
}
type YogaRawListQuery struct {
	UserId int `query:"userId"`
}

type YogaDeleteRawParam struct {
	Id int `params:"id"`
}
