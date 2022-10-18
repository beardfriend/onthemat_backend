package auth

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Auth interface {
	GenerateToken() (string, error)
	ParseToken(tokenString string, result jwt.Claims) error
}

var ErrInvalidToken = errors.New("invalid token")

type auth struct {
	options
	store Store
}

type options struct {
	signingMethod jwt.SigningMethod
	signKey       string
	claims        jwt.Claims
}

func NewAuth(options options) Auth {
	return &auth{
		options: options,
	}
}

func (a *auth) GenerateToken() (string, error) {
	signKey := []byte(a.options.signKey)

	token := jwt.NewWithClaims(a.options.signingMethod, a.claims)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *auth) ParseToken(tokenString string, result jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, result, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.options.signKey), nil
	})

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (a *auth) DestroyToken(ctx context.Context, tokenString string) error {
	var result jwt.Claims
	if err := a.ParseToken(tokenString, result); err != nil {
		return err
	}

	return a.store.Set(ctx, tokenString, 0)
}
