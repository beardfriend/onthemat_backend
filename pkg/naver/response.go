package naver

type GetTokenErrorBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GetTokenSuccessBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint32 `json:"expires_in"`
}

type GetUserInfo struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	NickName string `json:"nickname"`
	Name     string `json:"name"`
}
