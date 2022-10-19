package jwt

import (
	"fmt"
	"testing"
	"time"

	j "github.com/golang-jwt/jwt/v4"
)

func TestNewJwt(t *testing.T) {
	a := NewJwt(nil)
	fmt.Println(a.GetExpiredAt())
	s1, _ := a.GenerateToken()
	tokenInfo := &j.RegisteredClaims{}
	a.ParseToken(s1, tokenInfo)
	fmt.Println(tokenInfo)
}

func TestSetClaim(t *testing.T) {
	type customClaim struct {
		LoginType string
		j.RegisteredClaims
	}

	a := NewJwt(nil)
	custom := &customClaim{
		LoginType: "Admin",
		RegisteredClaims: j.RegisteredClaims{
			IssuedAt:  j.NewNumericDate(time.Now()),
			NotBefore: j.NewNumericDate(time.Now()),
			ExpiresAt: j.NewNumericDate(time.Now().Add(time.Duration(a.GetExpiredAt()) * time.Hour)),
			Issuer:    "onthemat",
		},
	}
	s1, _ := a.
		SetClaims(custom).
		GenerateToken()

	tokenInfo := &customClaim{}
	a.ParseToken(s1, tokenInfo)

	fmt.Println(tokenInfo)
}
