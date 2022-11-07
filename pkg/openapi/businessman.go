package openapi

import (
	"encoding/json"
	"fmt"
	"os"

	"onthemat/internal/app/config"

	"github.com/valyala/fasthttp"
)

type BusinessMan struct {
	Url    string
	Key    string
	client *fasthttp.Client
}

func NewBusinessMan(config *config.Config) *BusinessMan {
	return &BusinessMan{
		Url:    "https://api.odcloud.kr/api/nts-businessman/v1",
		Key:    config.APIKey.Businessman,
		client: &fasthttp.Client{},
	}
}

func (b *BusinessMan) GetStatus(businessNo string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(b.Url + "/status")

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.URI().QueryArgs().Add("serviceKey", b.Key)

	body := make(map[string]interface{})
	body["b_no"] = []string{businessNo}
	bodyBytes, _ := json.Marshal(body)
	req.SetBody(bodyBytes)

	resp := fasthttp.AcquireResponse()
	err := b.client.Do(req, resp)

	fasthttp.ReleaseRequest(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}
	return resp
}
