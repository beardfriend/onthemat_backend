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
