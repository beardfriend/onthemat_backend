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
		return c.Status(http.StatusUnauthorized).JSON(ex.NewBadRequestError("헤더를 확인해주세요."))
	}

	claim := &token.TokenClaim{}
	if err := m.tokensvc.ParseToken(access, claim); err != nil {

		if err.Error() == jwt.ErrExiredToken {
			return c.
				Status(http.StatusUnauthorized).
				JSON(ex.NewUnauthorizedError("토큰이 만료되었습니다. 재발급 해주세요."))
		}

		return c.Status(http.StatusBadRequest).
			JSON(ex.NewBadRequestError("토큰을 확인해주세요."))
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
			JSON(ex.NewUnauthorizedError("사업자 회원만 접근할 수 있습니다."))
	}

	return c.Next()
}

func (m *MiddleWare) OnlyTeacher(c *fiber.Ctx) error {
	userType := c.Context().UserValue("user_type").(string)

	if userType != "teacher" {
		return c.
			Status(http.StatusUnauthorized).
			JSON(ex.NewUnauthorizedError("선생님 회원만 접근할 수 있습니다."))
	}

	return c.Next()
}
