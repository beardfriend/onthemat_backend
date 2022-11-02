package http

import (
	"crypto/sha256"

	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	AuthUseCase usecase.AuthUseCase
	UserUseCase usecase.UserUseCase
}

func NewAuthHandler(
	authUseCase usecase.AuthUseCase,
	router fiber.Router,
) {
	handler := &authHandler{
		AuthUseCase: authUseCase,
	}
	g := router.Group("/auth")
	g.Get("/kakao", handler.Kakao)
	g.Get("/kakao/callback", handler.CallBackToken)
	g.Post("/signup", handler.SignUp)
	g.Post("/social/signup", handler.SocialSignUp)
}

func (h *authHandler) Kakao(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

func (h *authHandler) CallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	data := h.AuthUseCase.KakaoLogin(ctx, code)

	return c.JSON(data)
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.SignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := validatorx.ValidateStruct(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	body.Password = string(sha256.New().Sum([]byte(body.Password)))

	if err := h.AuthUseCase.SignUp(ctx, body); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}

func (h *authHandler) SocialSignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.SocialSignUpBody)
	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := h.AuthUseCase.SocialSignUp(ctx, body); err != nil {
		return c.SendStatus(500)
	}

	return c.SendStatus(200)
}
