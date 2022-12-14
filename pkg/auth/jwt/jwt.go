package jwt

import (
	"strings"

	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type JwtOption interface {
	WithSigningMethod(method jwtLib.SigningMethod) JwtOption
	WithSignKey(key string) JwtOption

	Init() Jwt
}

type Jwt interface {
	GenerateToken(claim jwtLib.Claims) (string, error)
	ParseToken(tokenString string, result jwtLib.Claims) error
}

const (
	ErrInvalidToken = "ErrInvalidToken"
	ErrExiredToken  = "ErrExpired"
)

type jwt struct {
	*option
}

func NewJwt() JwtOption {
	defaultOption := option{
		signingMethod: jwtLib.SigningMethodHS256,
		signKey:       "defaultSignKey",
	}
	return &jwt{
		option: &defaultOption,
	}
}

func (a *jwt) GenerateToken(claim jwtLib.Claims) (string, error) {
	signKey := []byte(a.option.signKey)

	token := jwtLib.NewWithClaims(a.option.signingMethod, claim)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *jwt) ParseToken(tokenString string, result jwtLib.Claims) error {
	token, err := jwtLib.ParseWithClaims(tokenString, result, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrInvalidToken)
		}
		return []byte(a.option.signKey), nil
	})

	if err != nil || !token.Valid {
		if strings.Contains(err.Error(), "expired") {
			return errors.New(ErrExiredToken)
		}
		return errors.New(ErrInvalidToken)
	}

	return nil
}
