package dto

type UserSignupBody struct{}

type AcademyNormalSignUpBody struct {
	NickName     string `json:"nickname"`
	PhoneNum     string `json:"phoneNum"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	BusinessCode string `json:"bussinessCode"`
	Address      string `json:"address"`
}

type LoginRequestQuery struct{}
