package http

import (
	"crypto/sha256"

	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	AuthUseCase usecase.AuthUseCase
	UserUseCase usecase.UserUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase, userUseCase usecase.UserUseCase, router fiber.Router) {
	handler := &authHandler{
		AuthUseCase: authUseCase,
		UserUseCase: userUseCase,
	}
	g := router.Group("/auth")
	g.Get("/kakao", handler.Kakao)
	g.Get("/kakao/callback", handler.CallBackToken)
	g.Get("/signup", handler.SignUp)
}

func (h *authHandler) Kakao(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

func (h *authHandler) CallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	code := c.Query("code")
	access, refresh := h.AuthUseCase.KakaoLogin(ctx, code)

	return c.JSON(fiber.Map{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *authHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	u, e := h.UserUseCase.GetMe(ctx, 1)

	if e != nil {
		panic(e)
	}
	resp := transport.NewUserMeResponse(u)
	return c.JSON(resp)
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	body := new(transport.SignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	body.Password = string(sha256.New().Sum([]byte(body.Password)))

	if err := h.AuthUseCase.SignUp(ctx, body); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func (h *authHandler) SocialSignUp(c *fiber.Ctx) error {
	ctx := c.Context()
	defer ctx.Done()

	body := new(transport.SocialSignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := h.AuthUseCase.SocialSignUp(ctx, body); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
