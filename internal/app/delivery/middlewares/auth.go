package middlewares

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/service/token"
	"onthemat/pkg/auth/jwt"

	"github.com/gofiber/fiber/v2"
)

func (m *middleWare) Auth(c *fiber.Ctx) error {
	ctx := c.Context()

	authorizationHeader := c.Request().Header.Peek("Authorization")

	access, err := m.authSvc.ExtractTokenFromHeader(string(authorizationHeader))
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrAuthorizationHeaderFormatUnavailable, "Bearer"))
	}

	claim := &token.TokenClaim{}
	if err := m.tokensvc.ParseToken(access, claim); err != nil {

		if err.Error() == jwt.ErrExiredToken {
			return c.
				Status(http.StatusUnauthorized).
				JSON(ex.NewHttpError(ex.ErrTokenExpired, nil))
		}

		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrTokenInvalid, nil))
	}

	ctx.SetUserValue("login_type", claim.LoginType)
	ctx.SetUserValue("user_type", claim.UserType)
	ctx.SetUserValue("user_id", claim.UserId)
	return c.Next()
}

func (m *middleWare) OnlySuperAdmin(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "superAdmin" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewHttpError(ex.ErrOnlySuperAdmin, nil))
	}

	return c.Next()
}

func (m *middleWare) OnlyAcademy(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "academy" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewHttpError(ex.ErrOnlyAcademy, nil))
	}

	return c.Next()
}

func (m *middleWare) OnlyTeacher(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "teacher" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewHttpError(ex.ErrOnlyTeacher, nil))
	}

	return c.Next()
}
