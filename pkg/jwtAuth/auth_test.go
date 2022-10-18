package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestCustomClaim(t *testing.T) {
	type AceessTokenPayload struct {
		LoginType string
		jwt.RegisteredClaims
	}

	jwtTime := time.Now().Add(time.Minute * 15)
	claim := &AceessTokenPayload{
		LoginType: "Admin",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(jwtTime),
			Issuer:    "onthemat",
		},
	}

	a := NewAuth(options{signingMethod: jwt.SigningMethodHS256, claims: claim, signKey: "asdsd"})
	token, err := a.GenerateToken()
	if err != nil {
		t.Error(err)
	}

	data := &AceessTokenPayload{}
	err = a.ParseToken(token, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data.LoginType)
}
