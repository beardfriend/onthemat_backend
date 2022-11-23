package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

type YogaGroupResponse struct {
	ID          int                  `json:"id"`
	Category    string               `json:"category"`
	CategoryEng string               `json:"categoryEng"`
	Description string               `json:"description"`
	CreatedAt   transport.TimeString `json:"createdAt"`
	UpdatedAt   transport.TimeString `json:"updatedAt"`
}

func NewYogaGroupsResponse(result []*ent.YogaGroup) []*YogaGroupResponse {
	response := make([]*YogaGroupResponse, 0)
	for _, v := range result {
		resp := new(YogaGroupResponse)
		copier.Copy(&resp, v)
		response = append(response, resp)
	}

	return response
}

type YogaResponse struct {
	ID          int                  `json:"id"`
	NameKor     string               `json:"nameKor"`
	NameEng     *string              `json:"nameEng"`
	Level       *int                 `json:"level"`
	Description *string              `json:"description"`
	CreatedAt   transport.TimeString `json:"createdAt"`
	UpdatedAt   transport.TimeString `json:"updatedAt"`
}

func NewYogaListResponse(result []*ent.Yoga) []*YogaResponse {
	response := make([]*YogaResponse, 0)
	for _, v := range result {
		resp := new(YogaResponse)
		copier.Copy(&resp, v)
		response = append(response, resp)
	}

	return response
}

type YogaRawResponse struct {
	ID        int                  `json:"id"`
	Name      string               `json:"name"`
	CreatedAt transport.TimeString `json:"createdAt"`
	UpdatedAt transport.TimeString `json:"updatedAt"`
}

func NewYogaRawListResponse(result []*ent.YogaRaw) []*YogaRawResponse {
	response := make([]*YogaRawResponse, 0)
	for _, v := range result {
		resp := new(YogaRawResponse)
		copier.Copy(&resp, v)
		response = append(response, resp)
	}

	return response
}
