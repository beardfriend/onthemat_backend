package transport

import (
	"time"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"

	"github.com/jinzhu/copier"
)

// ------------------- Request -------------------

type SignUpBody struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password"`
	NickName string `json:"nickname"`
}

type SocialSignUpBody struct {
	UserID    int    `json:"userId"`
	Email     string `json:"email"`
	TermAgree bool   `json:"termAgree"`
	NickName  string `json:"nickname"`
}

// ------------------- Response -------------------

type UserMeResponse struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Nickname    string    `json:"nickname"`
	CreatedAt   time.Time `json:"created_at"`
	SocialName  string    `json:"social_name"`
	SocialKey   int       `json:"social_key"`
	Type        user.Type `json:"type"`
	PhoneNum    string    `json:"phone_num"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func NewUserMeResponse(model *ent.User) *UserMeResponse {
	resp := new(UserMeResponse)
	copier.Copy(&resp, model)
	return resp
}
