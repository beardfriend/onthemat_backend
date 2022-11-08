package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type academyHandler struct {
	academyUsecase usecase.AcademyUsecase
}

func NewAcademyHandler(
	middleware *middlewares.MiddleWare,
	academyUsecase usecase.AcademyUsecase,
	router fiber.Router,
) {
	handler := &academyHandler{
		academyUsecase: academyUsecase,
	}

	g := router.Group("/academy")
	g.Post("", middleware.Auth, handler.Create)
}

func (h *academyHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(transport.AcademyCreateRequestBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.
			Status(http.StatusUnprocessableEntity).
			JSON(ex.NewUnprocessableEntityError("JSON을 입력해주세요."))
	}

	if err := h.academyUsecase.Create(ctx, reqBody, userId); err != nil {
		code, json := ex.ParseHttpError(err)
		return c.Status(code).JSON(json)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}
