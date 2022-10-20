package app

import (
	"time"

	"onthemat/internal/app/config"
	"onthemat/pkg/auth/jwt"

	j "github.com/golang-jwt/jwt/v4"
)

type Token interface {
	GenerateAccessToken()
}

type token struct {
	JWT   jwt.Jwt
	model *tokenModel
}

func NewToken(c *config.Config, jwt jwt.Jwt) *token {
	jwt.SetSignKey(c.JWT.SignKey)
	tokenModel := &tokenModel{
		RegisteredClaims: &j.RegisteredClaims{
			ExpiresAt: j.NewNumericDate(time.Now().Add(time.Duration(15) * time.Minute)),
		},
	}
	return &token{
		JWT:   jwt,
		model: tokenModel,
	}
}

func (t *token) SetLoginType(loginType string) *token {
	t.model.LoginType = loginType
	return t
}

func (t *token) SetExpiredAt(min int) *token {
	t.model.RegisteredClaims.ExpiresAt = j.NewNumericDate(time.Now().Add(time.Duration(min) * time.Minute))
	return t
}

func (t *token) SetUserID(id uint64) *token {
	t.model.UserID = id
	return t
}

type tokenModel struct {
	LoginType string
	UserID    uint64
	*j.RegisteredClaims
}

func (t *token) GenerateAccessToken(loginType string, id uint64) string {
	t.SetLoginType(loginType).SetUserID(id).SetExpiredAt(15)

	t.JWT.SetClaims(t.model)
	str, _ := t.JWT.GenerateToken()
	return str
}
