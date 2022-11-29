package request

// ------------------- SocialUrl -------------------

type SocialUrlParam struct {
	SocialName string `param:"socialName" validate:"required,oneof=kakao naver google"`
}

// ------------------- SocialCallback -------------------

type SocialCallbackParam struct {
	SocialName string `param:"socialName" validate:"required,oneof=kakao naver google"`
}
