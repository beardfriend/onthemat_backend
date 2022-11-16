package middlewares

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/service/token"
	"onthemat/pkg/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

func (m *MiddleWare) Auth(c *fiber.Ctx) error {
	ctx := c.Context()

	authorizationHeader := c.Request().Header.Peek("Authorization")

	access, err := m.authSvc.ExtractTokenFromHeader(string(authorizationHeader))
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewBadRequestError(ex.ErrAuthorizationHeaderFormatUnavailable, "Bearer"))
	}

	claim := &token.TokenClaim{}
	if err := m.tokensvc.ParseToken(access, claim); err != nil {

		if err.Error() == jwt.ErrExiredToken {
			return c.
				Status(http.StatusUnauthorized).
				JSON(ex.NewUnauthorizedError(ex.ErrTokenExpired, nil))
		}

		return c.Status(http.StatusBadRequest).
			JSON(ex.NewBadRequestError(ex.ErrTokenInvalid, nil))
	}

	ctx.SetUserValue("login_type", claim.LoginType)
	ctx.SetUserValue("user_type", claim.UserType)
	ctx.SetUserValue("user_id", claim.UserId)
	return c.Next()
}

func (m *MiddleWare) OnlyAcademy(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "academy" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewForbiddenError(ex.ErrOnlyAcademy, nil))
	}

	return c.Next()
}

func (m *MiddleWare) OnlyTeacher(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "teacher" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewForbiddenError(ex.ErrOnlyTeacher, nil))
	}

	return c.Next()
}
