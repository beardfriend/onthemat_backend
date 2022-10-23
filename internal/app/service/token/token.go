package token

import (
	"time"

	"onthemat/pkg/auth/jwt"

	jwtLib "github.com/golang-jwt/jwt/v4"
)

type TokenService interface {
	GenerateToken(uuid string, userId uint, loginType string, userType string, expired int) (string, error)
	ParseToken(tokenString string, result jwtLib.Claims) error
}

// ------------------- default -------------------

type tokenService struct {
	jwtPackage jwt.Jwt
}

const ErrMustSettingExpiredTime = "Please Set Exired Time (minute) "

func NewToken(jwtPackage jwt.Jwt) TokenService {
	return &tokenService{
		jwtPackage: jwtPackage,
	}
}

// ------------------- Model -------------------

type TokenClaim struct {
	Uuid      string
	UserId    uint
	LoginType string
	UserType  string
	jwtLib.RegisteredClaims
}

// ------------------- Service -------------------

func (t *tokenService) GenerateToken(uuid string, userId uint, loginType string, userType string, expired int) (string, error) {
	claim := TokenClaim{
		Uuid:      uuid,
		UserId:    userId,
		LoginType: loginType,
		UserType:  userType,
		RegisteredClaims: jwtLib.RegisteredClaims{
			Issuer:    "oneTheMat",
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(time.Duration(expired) * time.Minute)),
			NotBefore: jwtLib.NewNumericDate(time.Now().Add(time.Duration(expired) * time.Minute)),
		},
	}
	return t.jwtPackage.GenerateToken(claim)
}

func (t *tokenService) ParseToken(tokenString string, result jwtLib.Claims) error {
	return t.jwtPackage.ParseToken(tokenString, result)
}
