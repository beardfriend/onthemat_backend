package token

import (
	"testing"

	"onthemat/pkg/auth/jwt"

	"github.com/stretchr/testify/assert"
)

func TestLogic(t *testing.T) {
	jwt := jwt.NewJwt().WithSignKey("asd").Init()
	tokenModule := NewToken(jwt)

	to, _ := tokenModule.GenerateToken("uuid", 1, "kakao", "teacher", 10)
	cl := &TokenClaim{}
	tokenModule.ParseToken(to, cl)
	assert.Equal(t, cl.Uuid, "uuid")
	assert.Equal(t, cl.UserId, uint(1))
}
