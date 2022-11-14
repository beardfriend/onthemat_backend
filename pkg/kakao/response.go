package kakao

type GetTokenErrorBody struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type GetTokenSuccessBody struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             uint32 `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn uint64 `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

type GetUserInfoSuccessBody struct {
	Id           uint `json:"id"`
	KakaoAccount struct {
		Email *string `json:"email"`

		Profile struct {
			NickName string `json:"nickname"`
		} `json:"profile"`
	} `json:"kakao_account"`
}
