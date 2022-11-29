package request

// ------------------- SocialUrl -------------------

type AuthSocialUrlParam struct {
	SocialName string `param:"socialName" validate:"required,oneof=kakao naver google"`
}

// ------------------- SocialCallback -------------------

type AuthSocialCallbackParam struct {
	SocialName string `param:"socialName" validate:"required,oneof=kakao naver google"`
}

// ------------------- SignUp -------------------
type AuthSignUpBody struct {
	Email     string `json:"email" validate:"required,email,min=6,max=32"`
	Password  string `json:"password" validate:"required,min=10,max=20"`
	NickName  string `json:"nickname" validate:"required,min=2,max=10"`
	TermAgree bool   `json:"termAgree" validate:"required"`
}

// ------------------- SocialSignup -------------------

type AuthSocialSignUpBody struct {
	UserID int    `json:"userId" validate:"required"`
	Email  string `json:"email" validate:"required,email,min=6,max=32"`
}

// ------------------- Login -------------------

type AuthLoginBody struct {
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=10,max=20"`
}

// ------------------- CheckDuplicatedEmail -------------------

type AuthCheckDuplicatedEmailQuery struct {
	Email string `query:"email,required" validate:"required,email"`
}

// ------------------- SendTempPassword -------------------

type AuthTempPasswordQuery struct {
	Email string `query:"email,required" validate:"required,email"`
}

// ------------------- VerifyEmail -------------------
type AuthVerifyEmailQueries struct {
	Email    string `query:"email,required" validate:"email"`
	IssuedAt string `query:"isseudAt,required"`
	Key      string `query:"key,required"`
}
