package request

// ------------------- Update -------------------

type UserUpdateBody struct {
	Nickname string `json:"nickname"`
	PhoneNum string `json:"phone_num"`
}

type UserUpdateParam struct {
	Id int `params:"id" validate:"required"`
}
