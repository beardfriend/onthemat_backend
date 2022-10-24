package http

import (
	"context"
	"crypto/sha256"

	"onthemat/internal/app/dto"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/kakao"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	SignUp(c *fiber.Ctx) error
}

type authHandler struct {
	AuthUseCase usecase.AuthUseCase
	KakaoModule kakao.Kakao
}

func NewAuthHandler(kakao kakao.Kakao, authUseCase usecase.AuthUseCase) AuthHandler {
	return &authHandler{
		AuthUseCase: authUseCase,
	}
}

func (h *authHandler) KakaoLogin(c *fiber.Ctx) error {
	return c.Redirect(h.KakaoModule.Authorize())
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := context.Background()
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
