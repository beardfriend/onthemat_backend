package jwt

import (
	"fmt"
	"testing"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v4"
)

func TestInit(t *testing.T) {
	type customClaim struct {
		Name string
		jwtLib.RegisteredClaims
	}
	module := NewJwt().WithClaim(&customClaim{
		Name: "asdasd",
		RegisteredClaims: jwtLib.RegisteredClaims{
			Issuer:    "asd",
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Duration(100) * time.Minute)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}).Init()

	token, err := module.GenerateToken()
	if err != nil {
		t.Error(err)
	}
	result := &customClaim{}
	module.ParseToken(token, result)
	fmt.Println(result)
}
