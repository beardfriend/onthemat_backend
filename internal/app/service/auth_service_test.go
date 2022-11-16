package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	authService := NewAuthService(nil, nil, nil, nil)
	hashString := authService.HashPassword("password", "secret")
	hashString2 := authService.HashPassword("password", "secret")
	diffHahsString := authService.HashPassword("diffpassword", "secret")
	assert.NotEqual(t, diffHahsString, hashString)
	assert.Equal(t, hashString, hashString2)
}

func TestCheckEmailExpiredForVerify(t *testing.T) {
	as := NewAuthService(nil, nil, nil, nil)
	i := as.IsEmailForVerifyExpired("2022-11-16T13:04:05+09:00")
	assert.Equal(t, i, false)
	fmt.Println(time.Now().Format(time.RFC3339))
	d := as.IsEmailForVerifyExpired("2022-11-15T13:04:05+09:00")
	assert.Equal(t, d, true)
}
