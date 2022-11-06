package google

type GetTokenErrorBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type GetTokenSuccessBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
}

type GetUserInfo struct {
	Email string `json:"email"`
	Sub   uint   `json:"sub"`
}
