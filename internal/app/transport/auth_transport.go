package transport

import (
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
	UserID int    `json:"userId" validate:"required"`
	Email  string `json:"email" validate:"required,email,min=6,max=32"`
}

type CheckDuplicatedEmailQueries struct {
	Email string `query:"email,required" validate:"required,email"`
}

type SendTempPasswordQueries struct {
	Email string `query:"email,required" validate:"required,email"`
}

type VerifyEmailQueries struct {
	Email    string `query:"email,required" validate:"email"`
	IssuedAt string `query:"isseudAt,required"`
	Key      string `query:"key,required"`
}

// ------------------- Response -------------------

type UserMeResponse struct {
	ID          int        `json:"id"`
	Email       *string    `json:"email"`
	Nickname    *string    `json:"nickname"`
	SocialName  *string    `json:"social_name"`
	SocialKey   *string    `json:"social_key"`
	Type        *string    `json:"type"`
	PhoneNum    *string    `json:"phone_num"`
	CreatedAt   TimeString `json:"created_at"`
	LastLoginAt TimeString `json:"last_login_at"`
}

func NewUserMeResponse(model *ent.User) *UserMeResponse {
	resp := new(UserMeResponse)
	copier.Copy(&resp, model)
	resp.Type = model.Type.ToString()

	return resp
}
