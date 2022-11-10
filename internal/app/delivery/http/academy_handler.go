package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
	"onthemat/pkg/validatorx"

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
	g.Put("", middleware.Auth, handler.Update)
	g.Get("/:id", handler.Detail)
}

func (h *academyHandler) Detail(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(transport.AcademyDetailParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewBadRequestError("파라메터를 확인해주세요."))
	}
	if err := validatorx.ValidateStruct(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	academy, err := h.academyUsecase.Get(ctx, reqParam.Id)
	if err != nil {
		return utils.NewError(c, err)
	}

	response := transport.NewAcademyDetailResponse(academy)

	return c.Status(http.StatusCreated).JSON(ex.ResponseWithData{
		Code:    http.StatusCreated,
		Message: "",
		Result:  response,
	})
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
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *academyHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(transport.AcademyUpdateRequestBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.
			Status(http.StatusUnprocessableEntity).
			JSON(ex.NewUnprocessableEntityError("JSON을 입력해주세요."))
	}

	if err := h.academyUsecase.Update(ctx, reqBody, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}
