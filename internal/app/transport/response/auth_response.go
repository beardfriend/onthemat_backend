package response

import (
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	"github.com/jinzhu/copier"
)

type RefreshResponse struct {
	AccessToken          string               `json:"accessToken"`
	AccessTokenExpiredAt transport.TimeString `json:"accessTokenExpiredAt"`
}

func NewRefreshResponse(result *usecase.RefreshResult) *RefreshResponse {
	resp := new(RefreshResponse)
	copier.Copy(&resp, result)
	return resp
}
