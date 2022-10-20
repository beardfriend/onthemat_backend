package app

import (
	"fmt"
	"testing"

	"onthemat/internal/app/config"
	"onthemat/pkg/auth/jwt"
)

func TestNewAccessToken(t *testing.T) {
	c := config.NewConfig()
	if err := c.Load("../../configs"); err != nil {
		t.Error(err)
	}
	jwt := jwt.NewJwt()

	token := NewToken(c, jwt)
	access := token.GenerateAccessToken("kakao", 123)
	fmt.Println(access)
}
