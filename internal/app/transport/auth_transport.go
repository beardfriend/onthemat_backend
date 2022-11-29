package transport

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
