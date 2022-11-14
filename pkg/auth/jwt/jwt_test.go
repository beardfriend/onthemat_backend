package jwt

import (
	"testing"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

type customClaim struct {
	Name string
	jwtLib.RegisteredClaims
}

var customClm = customClaim{
	Name: "세훈",
	RegisteredClaims: jwtLib.RegisteredClaims{
		ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(50 * time.Second)),
	},
}

var normalClm = jwtLib.RegisteredClaims{
	ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(2 * time.Second)),
}

func TestJwtNormalConfig(t *testing.T) {
	assert := assert.New(t)

	t.Run("일반 생성", func(t *testing.T) {
		module := NewJwt().Init()

		token, err := module.GenerateToken(normalClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)
	})

	t.Run("커스텀 claim 생성", func(t *testing.T) {
		module := NewJwt().Init()
		token, err := module.GenerateToken(customClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)
	})

	t.Run("토큰 파싱 잘 되는지", func(t *testing.T) {
		module := NewJwt().Init()
		token, err := module.GenerateToken(customClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)

		var res customClaim
		if err := module.ParseToken(token, &res); err != nil {
			t.Error(err)
		}

		assert.Equal(res.Name, "세훈")
	})

	t.Run("만료기간이 지나면 에러 넘어오는지", func(t *testing.T) {
		module := NewJwt().WithSignKey("asdasd").Init()

		token, _ := module.GenerateToken(customClm)
		time.Sleep(3 * time.Second)

		var res customClaim
		if err := module.ParseToken(token, &res); err != nil {
			assert.Error(err)
		}
	})
}

func TestJwtCustomConfig(t *testing.T) {
	assert := assert.New(t)
	t.Run("일반 생성", func(t *testing.T) {
		module := NewJwt().WithSignKey("asdsdfkq").WithSigningMethod(jwtLib.SigningMethodHS512).Init()

		token, err := module.GenerateToken(normalClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)
	})

	t.Run("커스텀 claim 생성", func(t *testing.T) {
		module := NewJwt().WithSignKey("asdsdfkq").WithSigningMethod(jwtLib.SigningMethodHS512).Init()
		token, err := module.GenerateToken(customClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)
	})

	t.Run("토큰 파싱 잘 되는지", func(t *testing.T) {
		module := NewJwt().WithSignKey("asdsdfkq").WithSigningMethod(jwtLib.SigningMethodHS512).Init()
		token, err := module.GenerateToken(customClm)
		if err != nil {
			t.Error(err)
		}

		assert.GreaterOrEqual(len(token), 10)

		var res customClaim
		if err := module.ParseToken(token, &res); err != nil {
			t.Error(err)
		}

		assert.Equal(res.Name, "세훈")
	})

	t.Run("만료기간이 지나면 에러 넘어오는지", func(t *testing.T) {
		module := NewJwt().WithSignKey("asdsdfkq").WithSigningMethod(jwtLib.SigningMethodHS512).Init()

		token, _ := module.GenerateToken(customClm)
		time.Sleep(3 * time.Second)

		var res customClaim
		if err := module.ParseToken(token, &res); err != nil {
			assert.Error(err)
		}
	})
}
