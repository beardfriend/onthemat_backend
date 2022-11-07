package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"onthemat/pkg/email"
	"onthemat/pkg/ent"
	"onthemat/pkg/google"
	"onthemat/pkg/kakao"
	"onthemat/pkg/naver"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	ExtractTokenFromHeader(token string) (string, error)
	GetKakaoRedirectUrl() string
	GetGoogleRedirectUrl() string
	GetNaverRedirectUrl() string
	GetKakaoInfo(code string) (*kakao.GetUserInfoSuccessBody, error)
	GetGoogleInfo(code string) (*google.GetUserInfo, error)
	GetNaverInfo(code string) (*naver.GetUserInfo, error)

	GenerateRandomString() string
	GenerateRandomPassword() string
	SendEmailResetPassword(user *ent.User) error
	SendEmailVerifiedUser(email string, authKey string) error
}

type authService struct {
	kakao  *kakao.Kakao
	google *google.Google
	naver  *naver.Naver
	email  *email.Email
}

func NewAuthService(kakao *kakao.Kakao, google *google.Google, naver *naver.Naver, email *email.Email) AuthService {
	return &authService{
		kakao:  kakao,
		google: google,
		naver:  naver,
		email:  email,
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
	tokenResponse := a.kakao.GetToken(code)

	if tokenResponse.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(tokenResponse.Body(), body)

		return nil, errors.New(body.Error + body.ErrorCode)
	}

	tokenResponseBody := new(kakao.GetTokenSuccessBody)
	json.Unmarshal(tokenResponse.Body(), tokenResponseBody)

	infoResp := a.kakao.GetUserInfo(tokenResponseBody.AccessToken)

	if infoResp.StatusCode() != 200 {
		body := new(kakao.GetTokenErrorBody)
		json.Unmarshal(infoResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorCode)
	}

	infoRespBody := new(kakao.GetUserInfoSuccessBody)
	json.Unmarshal(infoResp.Body(), infoRespBody)

	return infoRespBody, nil
}

func (a *authService) GetGoogleInfo(code string) (*google.GetUserInfo, error) {
	tokenResp := a.google.GetToken(code)

	if tokenResp.StatusCode() != 200 {
		body := new(google.GetTokenErrorBody)
		json.Unmarshal(tokenResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorDescription)
	}

	tokenRespBody := new(google.GetTokenSuccessBody)
	json.Unmarshal(tokenResp.Body(), tokenRespBody)

	googleUserInfo := jwt.MapClaims{}
	jwt.ParseWithClaims(tokenRespBody.IdToken, &googleUserInfo, nil)

	infoRespBody := &google.GetUserInfo{
		Email: googleUserInfo["email"].(string),
		Sub:   googleUserInfo["sub"].(uint),
	}

	return infoRespBody, nil
}

func (a *authService) GetNaverInfo(code string) (*naver.GetUserInfo, error) {
	// naver
	tokenResp := a.naver.GetToken(code)
	if tokenResp.StatusCode() != 200 {
		body := new(naver.GetTokenErrorBody)
		json.Unmarshal(tokenResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorDescription)
	}

	tokenRespBody := new(naver.GetTokenSuccessBody)
	json.Unmarshal(tokenResp.Body(), tokenRespBody)

	// naver
	infoResp := a.naver.GetUserInfo(tokenRespBody.AccessToken)
	if infoResp.StatusCode() != 200 {
		body := new(naver.GetTokenErrorBody)
		json.Unmarshal(infoResp.Body(), body)

		return nil, errors.New(body.Error + body.ErrorDescription)
	}

	type res struct {
		Response naver.GetUserInfo `json:"response"`
	}
	response := new(res)

	json.Unmarshal(infoResp.Body(), &response)

	return &response.Response, nil
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

func (a *authService) GetNaverRedirectUrl() string {
	return a.naver.Authorize()
}

func (a *authService) SendEmailResetPassword(user *ent.User) error {
	subject := "임시 비밀번호 발급안내" + "!\n"
	body := "임시비밀번호는 " + user.TempPassword + " 입니다."
	msg := []byte(subject + "\n" + body)
	return a.email.Send([]string{user.Email}, msg)
}

func (a *authService) SendEmailVerifiedUser(email string, authKey string) error {
	href := fmt.Sprintf("http://localhost:3000/api/v1/auth/verify-email?key=%s&email=%s", authKey, email)
	subject := "Subject: Test email from Go!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
	<html>
		<body>
			<h1>이메일 인증입니다.</h1>
			<a href="%s">클릭</a>
		</body>
		</html>
		`, href)

	msg := []byte(subject + mime + body)
	return a.email.Send([]string{email}, msg)
}

func (a *authService) GenerateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 15
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (a *authService) GenerateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
