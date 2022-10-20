package jwt

import j "github.com/golang-jwt/jwt/v4"

type options struct {
	signingMethod j.SigningMethod
	signKey       string
	expired       int
	claims        j.Claims
}

func (a *jwt) SetClaims(claims j.Claims) Jwt {
	a.claims = claims
	return a
}

func (a *jwt) SetExpired(expired int) Jwt {
	a.expired = expired
	return a
}

func (a *jwt) SetSignKey(signKey string) Jwt {
	a.signKey = signKey
	return a
}

func (a *jwt) SetSigningMethod(signingMethod j.SigningMethod) Jwt {
	a.signingMethod = signingMethod
	return a
}

func (a *jwt) GetExpiredAt() int {
	return a.expired
}
