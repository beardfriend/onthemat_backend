package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/transport/response"
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
	g.Put("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	g.Patch("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	g.Delete("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	g.Get("/list", handler.List)
	g.Get("/:id", middleware.Auth, handler.Update)
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

func (h *recruitmentHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()
	academyId := ctx.UserValue("academy_id").(int)

	reqBody := new(request.RecruitmentUpdateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.RecruitmentUpdateParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	isUpdated, err := h.recruitmentUsecase.Update(ctx, reqBody, reqParam.Id, academyId)
	if err != nil {
		return utils.NewError(c, err)
	}

	httpCode := http.StatusOK
	if !isUpdated {
		httpCode = http.StatusCreated
	}

	return c.Status(httpCode).JSON(ex.Response{
		Code:    httpCode,
		Message: "",
	})
}

func (h *recruitmentHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	teacherId := ctx.UserValue("teacher_id").(int)

	reqBody := new(request.RecruitmentPatchBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.RecruitmentPatchParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	isUpdated, err := h.recruitmentUsecase.Patch(ctx, reqBody, reqParam.Id, teacherId)
	if err != nil {
		return utils.NewError(c, err)
	}
	httpCode := http.StatusOK
	if !isUpdated {
		httpCode = http.StatusCreated
	}

	return c.Status(httpCode).JSON(ex.Response{
		Code:    httpCode,
		Message: "",
	})
}

func (h *recruitmentHandler) List(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQueries := request.NewRecruitmentListQueries()

	if err := fiberx.QueryParser(c, reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	if err := h.Validator.ValidateStruct(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	recruits, pagination, err := h.recruitmentUsecase.List(ctx, reqQueries)
	resp := response.NewRecruitmentListResponse(recruits)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.ResponseWithPagination{
		Code:       http.StatusOK,
		Message:    "",
		Result:     resp,
		Pagination: pagination,
	})
}
