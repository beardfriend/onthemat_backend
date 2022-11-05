package http

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	ex "onthemat/internal/app/common"
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
	g.Get("/kakao/callback", handler.KakaoCallBackToken)
	g.Get("/google", handler.Google)
	g.Get("/google/callback", handler.GoogleCallBackToken)
	g.Post("/signup", handler.SignUp)
	g.Post("/social/signup", handler.SocialSignUp)
}

// Kakao godoc
// @Summary Kakao
// @Description Kakao
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/auth/kakao [get]
func (h *authHandler) Kakao(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

func (h *authHandler) KakaoCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	data := h.AuthUseCase.KakaoLogin(ctx, code)

	return c.JSON(data)
}

func (h *authHandler) Google(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.GoogleRedirectUrl(ctx))
}

func (h *authHandler) GoogleCallBackToken(c *fiber.Ctx) error {
	// ctx := c.Context()

	code := c.Query("code")
	fmt.Println(code)

	return c.JSON(code)
}

func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.SignUpBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ex.NewUnprocessableEntityError("JSON을 입력해주세요"))
	}

	if err := validatorx.ValidateStruct(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	body.Password = string(sha256.New().Sum([]byte(body.Password)))

	if err := h.AuthUseCase.SignUp(ctx, body); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(ex.NewInternalServerError(err))
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
