package google

import (
	"fmt"
	"os"
	"strings"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type Google struct {
	authUrl  string
	tokenUrl string
	client   *fasthttp.Client
	config   *config.Config
}

func NewGoogle(config *config.Config) *Google {
	return &Google{
		authUrl:  "https://accounts.google.com",
		tokenUrl: "https://oauth2.googleapis.com",
		client:   &fasthttp.Client{},
		config:   config,
	}
}

func (g *Google) Authorize() string {
	var scope []string

	scope = append(scope, "email")
	scope = append(scope, "profile")

	scopeString := strings.Join(scope, " ")
	redirectUrl := fmt.Sprintf("%s/o/oauth2/v2/auth?scope=%s&response_type=code&redirect_uri=%s&client_id=%s",
		g.authUrl,
		scopeString,
		g.config.Oauth.GoogleRedirect,
		g.config.Oauth.GoogleClientId,
	)
	return redirectUrl
}

func (g *Google) GetToken(code string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(g.tokenUrl + "/token")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.URI().QueryArgs().Add("code", code)
	req.URI().QueryArgs().Add("client_id", g.config.Oauth.GoogleClientId)
	req.URI().QueryArgs().Add("client_secret", g.config.Oauth.GoogleClientSecret)
	req.URI().QueryArgs().Add("redirect_uri", g.config.Oauth.GoogleRedirect)
	req.URI().QueryArgs().Add("grant_type", "authorization_code")

	req.Header.SetMethod(fasthttp.MethodPost)

	resp := fasthttp.AcquireResponse()
	err := g.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp
}
