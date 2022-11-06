package http

import (
	"crypto/sha256"
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
	g.Get("/naver", handler.Naver)
	g.Get("/naver/callback", handler.NaverCallBackToken)
	g.Post("/signup", handler.SignUp)
	g.Post("/login", handler.Login)
	g.Post("/social/signup", handler.SocialSignUp)
	g.Get("/reset-password", handler.SendResetPassword)
	g.Get("/check-email", handler.CheckDuplicatedEmail)
}

// Kakao godoc
// @Summary      카카오 로그인 URL
// @Description  카카오 로그인 or 회원가입 API
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      302
// @Failure      500
// @Router       /api/v1/auth/kakao [get]
func (h *authHandler) Kakao(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

func (h *authHandler) KakaoCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	data, err := h.AuthUseCase.SocialLogin(ctx, "kakao", code)
	if err != nil {
		panic(err)
	}

	return c.JSON(data)
}

func (h *authHandler) Google(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.GoogleRedirectUrl(ctx))
}

func (h *authHandler) GoogleCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	data, err := h.AuthUseCase.SocialLogin(ctx, "google", code)
	if err != nil {
		panic(err)
	}

	return c.JSON(data)
}

func (h *authHandler) Naver(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.NaverRedirectUrl(ctx))
}

func (h *authHandler) NaverCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	data, err := h.AuthUseCase.SocialLogin(ctx, "naver", code)
	if err != nil {
		panic(err)
	}

	return c.JSON(data)
}

// Kakao godoc
// @Summary      일반 회원가입
// @Description  이메일 비밀번호로 회원가입 하는 API
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      201
// @Failure      400 {object}
// @Failure      422 JSON을 입력해주세요
// @Failure      500
// @Router       /api/v1/auth/signup [post]
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

	return c.SendStatus(201)
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

func (h *authHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.LoginBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ex.NewUnprocessableEntityError("JSON을 입력해주세요"))
	}

	if err := validatorx.ValidateStruct(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	data, err := h.AuthUseCase.Login(ctx, body)
	if err != nil {
		panic(err)
	}

	return c.JSON(data)
}

func (h *authHandler) CheckDuplicatedEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(transport.CheckDuplicatedEmailQueries)

	if err := c.QueryParser(queries); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ex.NewUnprocessableEntityError("파라메터를 입력해주세요"))
	}

	if err := validatorx.ValidateStruct(queries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.CheckDuplicatedEmail(ctx, queries.Email); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(http.StatusBadRequest, err.Error(), nil))
	}

	return c.SendStatus(http.StatusOK)
}

func (h *authHandler) SendResetPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(transport.CheckDuplicatedEmailQueries)

	if err := c.QueryParser(queries); err != nil {
		return c.Status(http.StatusUnprocessableEntity).JSON(ex.NewUnprocessableEntityError("파라메터를 입력해주세요"))
	}

	if err := validatorx.ValidateStruct(queries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.SendEmailResetPassword(ctx, queries.Email); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(http.StatusBadRequest, err.Error(), nil))
	}

	return c.SendStatus(http.StatusOK)
}
