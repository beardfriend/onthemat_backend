package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
	"onthemat/pkg/fiberx"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type recruitmentHandler struct {
	recruitmentUsecase usecase.RecruitmentUsecase
	Validator          validatorx.Validator
	router             fiber.Router
}

func NewRecruitmentHandler(
	middleware middlewares.MiddleWare,
	recruitmentUsecase usecase.RecruitmentUsecase,
	Validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &recruitmentHandler{
		recruitmentUsecase: recruitmentUsecase,
		Validator:          Validator,
		router:             router,
	}

	g := router.Group("/recruitment")

	g.Post("", middleware.Auth, middleware.OnlyAcademy, handler.Create)
}

func (h *recruitmentHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	academyId := ctx.UserValue("academy_id").(int)

	reqBody := new(request.RecruitmentCreateBody)

	if err := fiberx.BodyParser(c, reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.recruitmentUsecase.Create(ctx, reqBody, academyId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}
