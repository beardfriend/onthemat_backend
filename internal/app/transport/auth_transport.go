package transport

import (
	"time"

	"onthemat/pkg/ent"

	"github.com/jinzhu/copier"
)

// ------------------- Request -------------------

type SignUpBody struct {
	Email     string `json:"email" validate:"required,email,min=6,max=32"`
	Password  string `json:"password" validate:"required,min=10,max=20"`
	NickName  string `json:"nickname" validate:"required,min=2,max=10"`
	TermAgree bool   `json:"termAgree" validate:"required"`
}

type LoginBody struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=10,max=20"`
}

type SocialSignUpBody struct {
	UserID   int    `json:"userId"`
	PhoneNum string `json:"phone_num"`
	Email    string `json:"email"`
}

type CheckDuplicatedEmailQueries struct {
	Email string `query:"email" validate:"required,email"`
}

type VerifyEmailQueries struct {
	Email string `json:"email" validate:"required,email"`
	Key   string `json:"key" validate:"required"`
}

// ------------------- Response -------------------

type UserMeResponse struct {
	ID          int       `json:"id"`
	Email       *string   `json:"email"`
	Nickname    *string   `json:"nickname"`
	SocialName  *string   `json:"social_name"`
	SocialKey   *string   `json:"social_key"`
	Type        *string   `json:"type"`
	PhoneNum    *string   `json:"phone_num"`
	CreatedAt   time.Time `json:"created_at"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func NewUserMeResponse(model *ent.User) *UserMeResponse {
	resp := new(UserMeResponse)
	copier.Copy(&resp, model)
	resp.Type = model.Type.ToString()
	return resp
}
