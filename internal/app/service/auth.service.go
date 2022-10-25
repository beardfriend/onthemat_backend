package service

import (
	"encoding/json"
	"errors"
	"strings"

	"onthemat/pkg/kakao"
)

type AuthService interface {
	ExtractTokenFromHeader(token string) (string, error)
	GetKakaoRedirectUrl() string
	GetKakaoID(code string) (string, error)
}

type authService struct {
	kakao *kakao.Kakao
}

func NewAuthService(kakao *kakao.Kakao) AuthService {
	return &authService{
		kakao: kakao,
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

func (a *authService) GetKakaoID(code string) (string, error) {
	// kakao
	resp := a.kakao.GetToken(code)
	if resp.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(resp.Body(), body)

		return "", errors.New(body.Error + body.ErrorCode)
	}

	// json
	body := new(kakao.GetTokenSuccessBody)
	json.Unmarshal(resp.Body(), body)

	// kakao
	respInfo := a.kakao.GetUserInfo(body.AccessToken)
	if respInfo.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(resp.Body(), body)

		return "", errors.New(body.Error + body.ErrorCode)
	}

	return body.AccessToken, nil
}

func (a *authService) GetKakaoRedirectUrl() string {
	// kakao
	resp := a.kakao.Authorize()
	r := resp.Header.Peek("Location")
	return string(r)
}
