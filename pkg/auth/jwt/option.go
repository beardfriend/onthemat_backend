package jwt

import (
	"reflect"

	jwtLib "github.com/golang-jwt/jwt/v4"
)

type option struct {
	signKey       string
	signingMethod jwtLib.SigningMethod
	claim         jwtLib.Claims
}

type optionForInit struct {
	option *option
}

func (a *jwt) WithSigningMethod(method jwtLib.SigningMethod) JwtOption {
	a.optionForInit.option.signingMethod = method
	return a
}

func (a *jwt) WithSignKey(key string) JwtOption {
	a.optionForInit.option.signKey = key
	return a
}

func (a *jwt) WithClaim(claim jwtLib.Claims) JwtOption {
	a.optionForInit.option.claim = claim
	return a
}

func (a *jwt) Init() Jwt {
	newOption := a.optionForInit.option
	defaultOption := a.option

	new := reflect.ValueOf(newOption).Elem()
	old := reflect.ValueOf(defaultOption).Elem()

	for i := 0; i < old.NumField(); i++ {
		if old.Field(i) != new.Field(i) {
			old.Set(new.Field(i))
		}
	}

	return a
}
