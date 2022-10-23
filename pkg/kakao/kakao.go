package kakao

import (
	"fmt"
	"os"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type kakao struct {
	AuthUrl string
	client  *fasthttp.Client
	config  *config.Config
}

func NewKakao(config *config.Config) *kakao {
	return &kakao{
		AuthUrl: "https://kauth.kakao.com/",
		client:  &fasthttp.Client{},
		config:  config,
	}
}

func (k *kakao) Authorize() string {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(k.AuthUrl)
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
	r := resp.Header.Peek("Location")

	return string(r)
}
