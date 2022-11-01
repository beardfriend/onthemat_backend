package middlewares

import (
	"net/http"

	"onthemat/internal/app/service/token"

	"github.com/gofiber/fiber/v2"
)

func (m *MiddleWare) Auth(c *fiber.Ctx) error {
	ctx := c.Context()

	authorizationHeader := c.Request().Header.Peek("Authorization")

	access, err := m.authSvc.ExtractTokenFromHeader(string(authorizationHeader))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON("에러")
	}

	claim := &token.TokenClaim{}
	if err := m.tokensvc.ParseToken(access, claim); err != nil {
		return c.Status(http.StatusUnauthorized).JSON("파싱")
	}

	ctx.SetUserValue("user_id", claim.UserId)
	return c.Next()
}
