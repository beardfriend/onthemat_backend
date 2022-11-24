package service

import (
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
