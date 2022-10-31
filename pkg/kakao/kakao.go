package kakao

import (
	"fmt"
	"os"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type Kakao struct {
	AuthUrl string
	ApiUrl  string
	client  *fasthttp.Client
	config  *config.Config
}

func NewKakao(config *config.Config) *Kakao {
	return &Kakao{
		AuthUrl: "https://kauth.kakao.com",
		ApiUrl:  "https://kapi.kakao.com",
		client:  &fasthttp.Client{},
		config:  config,
	}
}

func (k *Kakao) Authorize() *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(k.AuthUrl + "/oauth/authorize")
	req.Header.Add("Content-Type", "text/html")
	req.URI().QueryArgs().Add("client_id", k.config.Oauth.KaKaoClientId)
	req.URI().QueryArgs().Add("redirect_uri", k.config.Oauth.KaKaoRedirect)
	req.URI().QueryArgs().Add("response_type", "code")

	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	err := k.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp
}

func (k *Kakao) GetToken(code string) *fasthttp.Response {
	cnf := k.config.Oauth

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(k.AuthUrl + "/oauth/token")

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	body := fmt.Sprintf("grant_type=authorization_code&client_id=%s&redirect_uri=%s&code=%s", cnf.KaKaoClientId, cnf.KaKaoRedirect, code)
	req.SetBody([]byte(body))

	resp := fasthttp.AcquireResponse()
	err := k.client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp
}

func (k *Kakao) GetUserInfo(accessToken string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(k.ApiUrl + "/v2/user/me")

	authorizationValue := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Add("Authorization", authorizationValue)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	resp := fasthttp.AcquireResponse()
	err := k.client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp
}
