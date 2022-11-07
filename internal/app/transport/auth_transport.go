package transport

import (
	"time"

	"onthemat/pkg/ent"
	"onthemat/pkg/ent/user"

	"github.com/jinzhu/copier"
)

// ------------------- Request -------------------

type SignUpBody struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=10,max=20"`
	NickName string `json:"nickname" validate:"required,min=2,max=10"`
}

type LoginBody struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=10,max=20"`
}

type SocialSignUpBody struct {
	UserID    int    `json:"userId"`
	Email     string `json:"email"`
	TermAgree bool   `json:"termAgree"`
	NickName  string `json:"nickname"`
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
	Email       string    `json:"email"`
	Nickname    string    `json:"nickname"`
	CreatedAt   time.Time `json:"created_at"`
	SocialName  string    `json:"social_name"`
	SocialKey   string    `json:"social_key"`
	Type        user.Type `json:"type"`
	PhoneNum    string    `json:"phone_num"`
	LastLoginAt time.Time `json:"last_login_at"`
}

func NewUserMeResponse(model *ent.User) *UserMeResponse {
	resp := new(UserMeResponse)
	copier.Copy(&resp, model)
	return resp
}
