package jwt

import (
	"time"

	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type JwtOption interface {
	WithSigningMethod(method jwtLib.SigningMethod) JwtOption
	WithSignKey(key string) JwtOption
	WithClaim(claim jwtLib.Claims) JwtOption

	Init() Jwt
}

type Jwt interface {
	GenerateToken() (string, error)
	ParseToken(tokenString string, result jwtLib.Claims) error
}

var ErrInvalidToken = errors.New("ErrInvalidToken")

type jwt struct {
	*option
}

func NewJwt() JwtOption {
	defaultExpired := 2

	defaultOption := option{
		signingMethod: jwtLib.SigningMethodHS256,
		signKey:       "defaultSignKey",
		claim: &jwtLib.RegisteredClaims{
			Issuer:    "OnTheMat",
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Duration(defaultExpired) * time.Minute)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}
	return &jwt{
		option: &defaultOption,
	}
}

func (a *jwt) GenerateToken() (string, error) {
	signKey := []byte(a.option.signKey)

	token := jwtLib.NewWithClaims(a.option.signingMethod, a.option.claim)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	// a.optionDefaultSetting()

	return tokenString, nil
}

func (a *jwt) ParseToken(tokenString string, result jwtLib.Claims) error {
	token, err := jwtLib.ParseWithClaims(tokenString, result, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.option.signKey), nil
	})

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	return nil
}
