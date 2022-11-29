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
	// 요가 부분 수정
	g.Patch("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Patch)
	// 요가 삭제
	g.Delete("/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.Delete)
	// 그룹아이디 별 요가 리스트 조회
	g.Get("/list", middleware.Auth, handler.ListByGroupId)

	// 요가 그룹 생성
	g.Post("/group", middleware.Auth, middleware.OnlySuperAdmin, handler.CreateGroup)
	// 요가 그룹 수정
	g.Put("/group/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.UpdateGroup)
	// 요가 그룹 부분 수정
	g.Patch("/group/:id", middleware.Auth, middleware.OnlySuperAdmin, handler.PatchGroup)
	// 요가 그룹 멀티삭제
	g.Delete("/group/:ids", middleware.Auth, middleware.OnlySuperAdmin, handler.DeleteGroup)
	// 요가 그룹 리스트
	g.Get("/groups", middleware.Auth, handler.GetGroups)
}

// 요가 생성
/**
@api {post} /yoga 요가 생성
@apiName postYoga
@apiVersion 1.0.0
@apiGroup yoga
@apiDescription 요가 생성 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiBody {Number} yogaGroupId 요가 그룹의 아이디
@apiBody {String} nameKor 요가 이름 한국어
@apiBody {String} [nameEng] 요가 이름 영어
@apiBody {Number=1,2,3,4,5} [level] 요가 난이도
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidateError <code>400</code> code: 2000~
@apiError YogaGroupDoesNotExist <code>409</code> code: 4008
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaCreateBody)
	if err := fiberx.BodyParser(c, reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, err.Error()))
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

// 요가 수정
/**
@api {put} /yoga/:id 요가 수정
@apiName updateYoga
@apiVersion 1.0.0
@apiGroup yoga
@apiDescription 요가 수정 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 요가 아이디
@apiBody {Number} yogaGroupId 요가 그룹의 아이디
@apiBody {String} nameKor 요가 이름 한국어
@apiBody {String} [nameEng] 요가 이름 영어
@apiBody {Number=1,2,3,4,5} [level] 요가 난이도
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidateError <code>400</code> code: 2000~
@apiError YogaGroupDoesNotExist <code>409 Conflict</code> code: 4008
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaUpdateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.YogaUpdateParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.yogaUsecase.Update(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 요가 부분 수정
/**
@api {patch} /yoga/:id 요가 부분 수정
@apiName patchYoga
@apiVersion 1.0.0
@apiGroup yoga
@apiDescription 요가 부분 수정 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 요가 아이디
@apiBody {Number} [yogaGroupId] 요가 그룹의 아이디
@apiBody {String} [nameKor] 요가 이름 한국어
@apiBody {String} [nameEng] 요가 이름 영어
@apiBody {Number=1,2,3,4,5} [level] 요가 난이도
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidateError <code>400</code> code: 2000~
@apiError YogaGroupDoesNotExist <code>409</code> code: 4008
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaPatchBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.YogaPatchParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.yogaUsecase.Patch(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 요가 삭제
/**
@api {delet} /yoga/:id 요가 삭제
@apiName deleteYoga
@apiVersion 1.0.0
@apiGroup yoga
@apiDescription 요가 삭제 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 요가 아이디
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError InternalServerError <code>500</code> code: 500
*/
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

// 요가 리스트 조회
/**
@api {delet} /yoga/list 요가 리스트 조회
@apiName listYoga
@apiVersion 1.0.0
@apiGroup yoga
@apiDescription 요가 리스트 조회

@apiQuery {Number} groupId 요가 그룹 아이디
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object[]} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.nameKor 요가 이름 (한국어)
@apiSuccess {String} [result.nameEng] 요가 이름 (영어)
@apiSuccess {Number="1,2,3,4,5"} [result.level] 요가 난이도
@apiSuccess {String} [result.description] 요가에 대한 설명
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError QueryMissing <code>400</code> code: 3001
@apiError InternalServerError <code>500</code> code: 500
*/
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

// 요가 그룹 등록
/**
@api {post} /yoga/group 요가 그룹 생성
@apiName postYogaGroup
@apiVersion 1.0.0
@apiGroup yogaGroup
@apiDescription 요가 생성 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiBody {String} category 요가 카테고리 (한국어)
@apiBody {String} categoryEng 요가 카테고리 (영어)
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidateErrors <code>400</code> code: 2000~
@apiError ErrYogaGroupAlreadyExist <code>409</code> code: 4008
@apiError InternalServerError <code>500</code> code: 500
*/
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

// 요가 그룹 수정
/**
@api {put} /yoga/group/:id 요가 그룹 수정
@apiName putYogaGroup
@apiVersion 1.0.0
@apiGroup yogaGroup
@apiDescription 요가 그룹 수정 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 요가 그룹 아이디
@apiBody {String} category 요가 카테고리 (한국어)
@apiBody {String} categoryEng 요가 카테고리 (영어)
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidateErrors <code>400</code> code: 2000~
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) UpdateGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaGroupUpdateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.YogaUpdateParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.yogaUsecase.UpdateGroup(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 요가 그룹 부분 수정
/**
@api {patch} /yoga/group/:id 요가 그룹 부분 수정
@apiName patchYogaGroup
@apiVersion 1.0.0
@apiGroup yogaGroup
@apiDescription 요가 그룹 부분 수정 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 요가 그룹 아이디
@apiBody {String} [category] 요가 카테고리 (한국어)
@apiBody {String} [categoryEng] 요가 카테고리 (영어)
@apiBody {String} [description] 요가에 대한 설명
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidateErrors <code>400</code> code: 2000~
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) PatchGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	reqBody := new(request.YogaGroupPatchBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.YogaGroupPatchParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	if err := h.yogaUsecase.PatchGroup(ctx, reqBody, reqParam.Id); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 요가 그룹 삭제
/**
@api {delete} /yoga/group/:ids 요가 그룹 삭제
@apiName deleteYogaGroup
@apiVersion 1.0.0
@apiGroup yogaGroup
@apiDescription 요가 그룹 삭제 (어드민 권한만)
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} ids 요가 그룹 아이디 ,로 멀티삭제
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiError ParamsMissing <code>400</code> code: 3002
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *yogaHandler) DeleteGroup(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.YogaGroupDeleteParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	_, err := h.yogaUsecase.DeleteGroup(ctx, reqParam.Ids)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.Response{
		Code:    http.StatusOK,
		Message: "",
	})
}

// 요가 그룹 리스트 조회
/**
@api {delet} /yoga/groups 요가 그룹 리스트 조회
@apiName listYogaGroup
@apiVersion 1.0.0
@apiGroup yogaGroup
@apiDescription 요가 그룹 리스트 조회
@apiQuery {Number} [pageNo=1] 페이지 번호
@apiQuery {Number} [pageSize=10] 페이지 당 문서 개수
@apiQuery {String} [category] 검색할 카테고리 이름
@apiQuery {String="DESC,ASC"} [orderType=DESC] 정렬 타입
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object[]} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.category 요가 그룹 이름 (한국어)
@apiSuccess {String} result.categoryEng 요가 그룹 이름 (영어)
@apiSuccess {String} [result.description] 요가 그룹에 대한 설명
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError QueryMissing <code>400</code> code: 3001
@apiError InternalServerError <code>500</code> code: 500
*/
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
