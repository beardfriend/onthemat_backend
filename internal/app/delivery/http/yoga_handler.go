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

	// 요가 등록
	g.Post("/", middleware.Auth, middleware.OnlySuperAdmin, handler.Create)
	// 요가 수정
	g.Put("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Update)
	// 요가 삭제
	g.Delete("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Delete)
	// 그룹아이디 별 요가 리스트 조회
	g.Get("/list", middleware.Auth, handler.ListByGroupId)

	// 요가 그룹 생성
	g.Post("/group", middleware.Auth, middleware.OnlySuperAdmin, handler.CreateGroup)
	// 요가 그룹 수정
	g.Put("/group/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.UpdateGroup)
	// 요가 그룹 멀티삭제
	g.Delete("/groups", middleware.Auth, middleware.OnlySuperAdmin, handler.DeleteGroups)
	// 요가 그룹 리스트
	g.Get("/group/list", middleware.Auth, handler.GetGroups)
	g.Patch("/test", handler.Patch)
}

func (h *yogaHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	reqBody := new(request.YogaPatcheBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.yogaUsecase.Patch(ctx, reqBody, 2); err != nil {
		return utils.NewError(c, err)
	}
	return c.Status(http.StatusCreated).JSON(reqBody)
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
