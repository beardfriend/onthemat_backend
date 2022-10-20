package jwt

import (
	"time"

	j "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Jwt interface {
	// Get Set
	GetExpiredAt() int

	SetSignKey(signKey string) *jwt
	SetClaims(claims j.Claims) *jwt

	// Function
	GenerateToken() (string, error)
	ParseToken(tokenString string, result j.Claims) error
}

var ErrInvalidToken = errors.New("ErrInvalidToken")

type jwt struct {
	options
}

func NewJwt() *jwt {
	expired := 1
	return &jwt{
		options: options{
			signingMethod: j.SigningMethodHS256,
			signKey:       "asd",
			expired:       expired,
			claims: j.RegisteredClaims{
				Issuer:    "ontheMat",
				ExpiresAt: j.NewNumericDate(time.Now().Add(time.Duration(expired) * time.Hour)),
				IssuedAt:  j.NewNumericDate(time.Now()),
				NotBefore: j.NewNumericDate(time.Now()),
				Subject:   "normal",
			},
		},
	}
}

func (a *jwt) GenerateToken() (string, error) {
	signKey := []byte(a.options.signKey)

	token := j.NewWithClaims(a.options.signingMethod, a.claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *jwt) ParseToken(tokenString string, result j.Claims) error {
	token, err := j.ParseWithClaims(tokenString, result, func(token *j.Token) (interface{}, error) {
		if _, ok := token.Method.(*j.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.options.signKey), nil
	})

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	return nil
}
