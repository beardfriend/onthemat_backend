package request

import (
	"onthemat/internal/app/utils"
)

// ------------------- Body -------------------
type YogaGroupCreateBody struct {
	Category    string  `json:"category" validate:"required"`
	CategoryEng string  `json:"categoryEng"`
	Description *string `json:"description"`
}

type YogaGroupUpdateBody struct {
	Category    string  `json:"category" validate:"required"`
	CategoryEng string  `json:"categoryEng"`
	Description *string `json:"description"`
}

type YogaGroupPatchBody struct {
	Category    *string `json:"category"`
	CategoryEng *string `json:"categoryEng"`
	Description *string `json:"description"`
}

type YogaCreateBody struct {
	YogaGroupId int     `json:"yogaGroupId"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"description"`
}

type YogaUpdateBody struct {
	YogaGroupId int     `json:"yogaGroupId" validate:"required"`
	NameKor     string  `json:"nameKor"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level"`
	Description *string `json:"description"`
}

type YogaPatchBody struct {
	YogaGroupId *int    `json:"yogaGroupId,omitempty" validate:"required"`
	NameKor     *string `json:"nameKor,omitempty"`
	NameEng     *string `json:"nameEng"`
	Level       *int    `json:"level,omitempty"`
	Description *string `json:"description,omitempty"`
}

type YogaGroupsDeleteBody struct {
	Ids []int `json:"ids" validate:"required"`
}

type YogaRawCreateBody struct {
	YogaName string `json:"name" validate:"required"`
}

// ------------------- Query String -------------------

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

type YogaListQuery struct {
	GroupId int `query:"groupId"`
}

type YogaRawListQuery struct {
	UserId int `query:"userId"`
}

// ------------------- Params -------------------

type YogaUpdateParam struct {
	Id int `params:"id"`
}

type YogaDeleteParam struct {
	Id int `params:"id"`
}

type YogaDeleteRawParam struct {
	Id int `params:"id"`
}
