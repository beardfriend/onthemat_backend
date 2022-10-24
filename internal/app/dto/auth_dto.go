package dto

type UserSignupBody struct{}

type SignUpBody struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	NickName     string `json:"nickname"`
}

type LoginRequestQuery struct{}
