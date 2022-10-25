package http

import (
	"crypto/sha256"

	"onthemat/internal/app/dto"
	"onthemat/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	AuthUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, router fiber.Router) {
	handler := &authHandler{
		AuthUseCase: authUseCase,
	}
	g := router.Group("/auth")
	g.Get("/kakao", handler.KakaoLogin)
	g.Get("/kakao/callback", handler.KakaoLoginCallBack)
}

func (h *authHandler) KakaoLogin(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	// 리프레쉬 있으면

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

func (h *authHandler) KakaoLoginCallBack(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	code := c.Query("code")
	h.AuthUseCase.KakaoLogin(ctx, code)

	return c.SendStatus(200)
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	body := new(dto.SignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	//
	// validation

	// password Hash

	body.Password = string(sha256.New().Sum([]byte(body.Password)))

	if err := h.AuthUseCase.SignUp(ctx, body); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
