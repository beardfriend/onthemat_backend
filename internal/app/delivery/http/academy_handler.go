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

type academyHandler struct {
	academyUsecase usecase.AcademyUsecase
	Validator      validatorx.Validator
}

func NewAcademyHandler(
	middleware middlewares.MiddleWare,
	academyUsecase usecase.AcademyUsecase,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &academyHandler{
		academyUsecase: academyUsecase,
		Validator:      validator,
	}

	g := router.Group("/academy")
	// 학원 등록
	g.Post("", middleware.Auth, handler.Create)
	// 학원 수정
	g.Put("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	// 학원 일부 수정
	g.Patch("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Patch)
	// 학원 리스트
	g.Get("/list", handler.List)
	// 학원 상세조회
	g.Get("/:id", handler.Detail)
}

func (h *academyHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(request.AcademyCreateBody)

	if err := fiberx.BodyParser(c, reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
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

	reqBody := new(request.AcademyUpdateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.AcademyUpdateParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	// 소유권체크
	if userId != reqParam.Id {
		return c.Status(403).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	if err := h.academyUsecase.Update(ctx, reqBody, reqParam.Id, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *academyHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)

	reqBody := new(request.AcademyPatchBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.AcademyPatchParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	// 소유권체크
	if userId != reqParam.Id {
		return c.Status(403).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	if err := h.academyUsecase.Patch(ctx, reqBody, reqParam.Id, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *academyHandler) Detail(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.AcademyDetailParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}
	if err := h.Validator.ValidateStruct(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	academy, err := h.academyUsecase.Get(ctx, reqParam.Id)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewAcademyDetailResponse(academy)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  resp,
	})
}

func (h *academyHandler) List(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQueries := request.NewAcademyListQueries()

	if err := c.QueryParser(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	academies, pagination, err := h.academyUsecase.List(ctx, reqQueries)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewAcademyListResponse(academies)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithPagination{
		Code:       http.StatusOK,
		Message:    "",
		Result:     resp,
		Pagination: pagination,
	})
}
