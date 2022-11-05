package google

import (
	"fmt"
	"strings"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type Google struct {
	authUrl string
	client  *fasthttp.Client
	config  *config.Config
}

func NewGoogle(config *config.Config) *Google {
	return &Google{
		authUrl: "https://accounts.google.com",
		client:  &fasthttp.Client{},
		config:  config,
	}
}

func (g *Google) Authorize() string {
	var scope []string

	scope = append(scope, "email")

	scopeString := strings.Join(scope, ", ")
	redirectUrl := fmt.Sprintf("%s/o/oauth2/v2/auth?scope=%s&response_type=code&redirect_uri=%s&client_id=%s",
		g.authUrl,
		scopeString,
		g.config.Oauth.GoogleRedirect,
		g.config.Oauth.GoogleClientId,
	)
	return redirectUrl
}
