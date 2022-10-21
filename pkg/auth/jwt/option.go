package jwt

import (
	jwtLib "github.com/golang-jwt/jwt/v4"
)

type option struct {
	signKey       string
	signingMethod jwtLib.SigningMethod
	claim         jwtLib.Claims
}

func (a *jwt) WithSigningMethod(method jwtLib.SigningMethod) JwtOption {
	a.option.signingMethod = method
	return a
}

func (a *jwt) WithSignKey(key string) JwtOption {
	a.option.signKey = key
	return a
}

func (a *jwt) WithClaim(claim jwtLib.Claims) JwtOption {
	a.option.claim = claim
	return a
}

func (a *jwt) Init() Jwt {
	return a
}
