package service

import (
	"fmt"
	"testing"

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

func TestIsExpiredEmailForVerify(t *testing.T) {
	as := NewAuthService(nil, nil, nil, nil)
	b := as.IsExpiredEmailForVerify("2022-12-26T11:18:26+09:00")
	fmt.Println(b)
}
