package jwt

import (
	"time"

	j "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Jwt interface {
	// Get Set
	GetExpiredAt() int
	SetClaims(claims j.Claims) Jwt

	// Function
	GenerateToken() (string, error)
	ParseToken(tokenString string, result j.Claims) error
}

var ErrInvalidToken = errors.New("ErrInvalidToken")

type jwt struct {
	options
}

type options struct {
	signingMethod j.SigningMethod
	signKey       string
	expiredAt     int
	claims        j.Claims
}

func NewAuth(opt *options) Jwt {
	expired := 24
	if opt == nil {
		return &jwt{
			options: options{
				signingMethod: j.SigningMethodHS256,
				signKey:       "asd",
				expiredAt:     expired,
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
	return &jwt{
		options: *opt,
	}
}

func (a *jwt) SetClaims(claims j.Claims) Jwt {
	a.claims = claims
	return a
}

func (a *jwt) GetExpiredAt() int {
	return a.expiredAt
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
