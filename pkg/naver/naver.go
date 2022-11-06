package naver

import (
	"fmt"
	"os"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type Naver struct {
	authUrl    string
	openAPIUrl string
	client     *fasthttp.Client
	config     *config.Config
}

func NewNaver(config *config.Config) *Naver {
	return &Naver{
		authUrl:    "https://nid.naver.com/oauth2.0",
		openAPIUrl: "https://openapi.naver.com",
		client:     &fasthttp.Client{},
		config:     config,
	}
}

func (n *Naver) Authorize() string {
	redirectUrl := fmt.Sprintf("%s/authorize?response_type=code&redirect_uri=%s&client_id=%s&state=%s",
		n.authUrl,
		n.config.Oauth.NaverRedirect,
		n.config.Oauth.NaverClientId,
		"randomText",
	)
	return redirectUrl
}

func (n *Naver) GetToken(code string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(n.authUrl + "/token")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URI().QueryArgs().Add("grant_type", "authorization_code")
	req.URI().QueryArgs().Add("client_id", n.config.Oauth.NaverClientId)
	req.URI().QueryArgs().Add("client_secret", n.config.Oauth.NaverClientSecret)
	req.URI().QueryArgs().Add("code", code)
	req.URI().QueryArgs().Add("state", "randomText")

	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	err := n.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp
}

func (k *Naver) GetUserInfo(accessToken string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(k.openAPIUrl + "/v1/nid/me")

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
