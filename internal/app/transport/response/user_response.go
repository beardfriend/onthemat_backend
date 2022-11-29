package response

import (
	"onthemat/internal/app/transport"
	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

// ------------------- Get -------------------

type UserMeResponse struct {
	ID          int                  `json:"id"`
	Email       *string              `json:"email"`
	Nickname    *string              `json:"nickname"`
	SocialName  *string              `json:"social_name"`
	SocialKey   *string              `json:"social_key"`
	Type        *string              `json:"type"`
	PhoneNum    *string              `json:"phone_num"`
	CreatedAt   transport.TimeString `json:"created_at"`
	LastLoginAt transport.TimeString `json:"last_login_at"`
}

func NewUserMeResponse(model *ent.User) *UserMeResponse {
	resp := new(UserMeResponse)
	copier.Copy(&resp, model)

	resp.SocialName = model.SocialName.ToString()
	resp.Type = model.Type.ToString()

	return resp
}
