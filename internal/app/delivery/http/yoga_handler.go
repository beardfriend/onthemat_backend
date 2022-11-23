package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/transport/response"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type yogaHandler struct {
	yogaUsecase usecase.YogaUsecase
	middleware  middlewares.MiddleWare
	validator   validatorx.Validator
	router      fiber.Router
}

func NewYogaHandler(
	yogaUsecase usecase.YogaUsecase,
	middleware middlewares.MiddleWare,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &yogaHandler{
		yogaUsecase: yogaUsecase,
		middleware:  middleware,
		validator:   validator,
		router:      router,
	}
	g := router.Group("/yoga")

	g.Post("/", middleware.Auth, middleware.OnlySuperAdmin, handler.Create)
	g.Put("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Update)
	g.Delete("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Delete)
	g.Get("/list", middleware.Auth, handler.ListByGroupId)

	g.Post("/group", middleware.Auth, middleware.OnlySuperAdmin, handler.CreateGroup)
	g.Put("/group/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.UpdateGroup)
	g.Delete("/groups", middleware.Auth, middleware.OnlySuperAdmin, handler.DeleteGroups)
	g.Get("/group/list", middleware.Auth, handler.GetGroups)

	g.Post("/raw", middleware.Auth, handler.CreateRaw)
	g.Delete("/raw/:id", middleware.Auth, handler.DeleteRaw)
	g.Get("/raw/list", handler.ListRawByUserId)
}

func (h *yogaHandler) CreateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaGroupCreateBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.yogaUsecase.CreateGroup(ctx, reqBody); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *yogaHandler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaGroupUpdateBody)
	reqParam := new(request.YogaUpdateParam)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.yogaUsecase.UpdateGroup(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *yogaHandler) DeleteGroups(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaGroupsDeleteBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	_, err := h.yogaUsecase.DeleteGroup(ctx, reqBody.Ids)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.Response{
		Code:    http.StatusOK,
		Message: "",
	})
}

func (h *yogaHandler) GetGroups(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQueries := request.NewYogaGroupListQueries()
	if err := c.QueryParser(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	result, paginationInfo, err := h.yogaUsecase.GroupList(ctx, reqQueries)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewYogaGroupsResponse(result)
	return c.Status(http.StatusOK).JSON(ex.ResponseWithPagination{
		Code:       http.StatusOK,
		Message:    "",
		Result:     resp,
		Pagination: paginationInfo,
	})
}

func (h *yogaHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaCreateBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.yogaUsecase.Create(ctx, reqBody); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *yogaHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaUpdateBody)
	reqParam := new(request.YogaUpdateParam)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.yogaUsecase.Update(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *yogaHandler) ListByGroupId(c *fiber.Ctx) error {
	ctx := c.Context()
	reqQuery := new(request.YogaListQuery)

	if err := c.QueryParser(reqQuery); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}
	data, err := h.yogaUsecase.List(ctx, reqQuery.GroupId)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewYogaListResponse(data)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  resp,
	})
}

func (h *yogaHandler) Delete(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.YogaDeleteParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}

	if err := h.yogaUsecase.Delete(ctx, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.Response{
		Code:    http.StatusOK,
		Message: "",
	})
}

func (h *yogaHandler) CreateRaw(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := c.Context().UserValue("user_id").(int)

	reqBody := new(request.YogaRawCreateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.yogaUsecase.CreateRaw(ctx, reqBody, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

func (h *yogaHandler) ListRawByUserId(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQuery := new(request.YogaRawListQuery)

	if err := c.QueryParser(reqQuery); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	data, err := h.yogaUsecase.ListRawByUserId(ctx, reqQuery.UserId)
	if err != nil {
		return utils.NewError(c, err)
	}
	resp := response.NewYogaRawListResponse(data)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  resp,
	})
}

func (h *yogaHandler) DeleteRaw(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := c.Context().UserValue("user_id").(int)

	reqParam := new(request.YogaDeleteRawParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}

	if _, err := h.yogaUsecase.DeleteRaw(ctx, reqParam.Id, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.Response{
		Code:    http.StatusOK,
		Message: "",
	})
}
