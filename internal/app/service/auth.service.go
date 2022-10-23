package service

import (
	"errors"
	"strings"
)

type AuthService interface {
	ExtractTokenFromHeader(token string) (string, error)
}

type authService struct{}

var ErrNotBearerToken = "Token unavailable"

func (a *authService) ExtractTokenFromHeader(token string) (string, error) {
	splitedToken := strings.Split(token, " ")
	if splitedToken[0] != "Bearer" {
		return "", errors.New(ErrNotBearerToken)
	}

	return splitedToken[1], nil
}
