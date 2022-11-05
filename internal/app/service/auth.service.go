package service

import (
	"encoding/json"
	"errors"
	"strings"

	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
)

type AuthService interface {
	ExtractTokenFromHeader(token string) (string, error)
	GetKakaoRedirectUrl() string
	GetGoogleRedirectUrl() string
	GetKakaoInfo(code string) (*kakao.GetUserInfoSuccessBody, error)
}

type authService struct {
	kakao  *kakao.Kakao
	google *google.Google
}

func NewAuthService(kakao *kakao.Kakao, google *google.Google) AuthService {
	return &authService{
		kakao:  kakao,
		google: google,
	}
}

var ErrNotBearerToken = "Token unavailable"

func (a *authService) ExtractTokenFromHeader(token string) (string, error) {
	splitedToken := strings.Split(token, " ")
	if splitedToken[0] != "Bearer" {
		return "", errors.New(ErrNotBearerToken)
	}

	return splitedToken[1], nil
}

func (a *authService) GetKakaoInfo(code string) (*kakao.GetUserInfoSuccessBody, error) {
	// kakao
	tokenResp := a.kakao.GetToken(code)
	if tokenResp.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(tokenResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorCode)
	}

	tokenRespBody := new(kakao.GetTokenSuccessBody)
	json.Unmarshal(tokenResp.Body(), tokenRespBody)

	// kakao
	infoResp := a.kakao.GetUserInfo(tokenRespBody.AccessToken)
	if infoResp.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(infoResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorCode)
	}

	infoRespBody := new(kakao.GetUserInfoSuccessBody)
	json.Unmarshal(infoResp.Body(), infoRespBody)
	return infoRespBody, nil
}

func (a *authService) GetKakaoRedirectUrl() string {
	// kakao
	resp := a.kakao.Authorize()
	r := resp.Header.Peek("Location")
	return string(r)
}

func (a *authService) GetGoogleRedirectUrl() string {
	return a.google.Authorize()
}
