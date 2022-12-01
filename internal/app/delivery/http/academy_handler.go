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
	// 학원 생성
	g.Post("", middleware.Auth, handler.Create)
	// 학원 수정
	g.Put("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	// 학원 일부 수정
	g.Patch("/:id", middleware.Auth, middleware.OnlyAcademy, handler.Patch)
	// 학원 리스트
	g.Get("/list", handler.List)
	// 학원 상세조회
	g.Get("/:id", handler.Get)
}

// 학원 생성
/**
@api {post} /academy 학원 생성
@apiName postAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 생성
@apiHeader Authorization accessToken (Bearer)
@apiBody {Object} info
@apiBody {String} info.sigunguId 시군구 id
@apiBody {String} info.businessCode 사업자 번호
@apiBody {String} info.name 학원 이름
@apiBody {String} info.callNumber 연락가능한 번호
@apiBody {String} info.addressRoad 도로명 주소
@apiBody {String} [info.addressDetail] 상세 주소
@apiBody {int[]} [yogaIds] 요가 아이디
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError BusinessCodeInvalid <code>400</code> code: 3008
@apiError UserNotFound <code>404</code> code: 5001
@apiError UserTypeAlreadyRegisted <code>409</code> code: 4003
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
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

// 학원 수정
/**
@api {put} /academy/:id 학원 수정
@apiName putAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 학원 아이디
@apiBody {Object} info
@apiBody {String} info.sigunguId 시군구 id
@apiBody {String} info.name 학원 이름
@apiBody {String} info.callNumber 연락가능한 번호
@apiBody {String} info.addressRoad 도로명 주소
@apiBody {String} [info.addressDetail] 상세 주소
@apiBody {int[]} [yogaIds] 요가 아이디
@apiSuccess (200 or 201) {Number} code 200 or 201
@apiSuccess (200 or 201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidationError <code>400</code> code: 2xxx
@apiError OnlyAcademy <code>403</code> code: 6004
@apiError OnlyOwnUser <code>403</code> code: 6007
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *academyHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)
	academy_id := ctx.UserValue("academy_id").(int)

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
	if academy_id != reqParam.Id {
		return c.Status(http.StatusForbidden).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	isUpdated, err := h.academyUsecase.Put(ctx, reqBody, reqParam.Id, userId)
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

// 학원 부분 수정
/**
@api {patch} /academy/:id 학원 부분 수정
@apiName patchAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 부분 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 학원 아이디
@apiBody {Object} [info]
@apiBody {String} [info.sigunguId] 시군구 id
@apiBody {String} [info.name] 학원 이름
@apiBody {String} [info.callNumber] 연락가능한 번호
@apiBody {String} [info.addressRoad] 도로명 주소
@apiBody {String} [info.addressDetail] 상세 주소
@apiBody {int[]} [yogaIds] 요가 아이디
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidationError <code>400</code> code: 2xxx
@apiError OnlyAcademy <code>403</code> code: 6004
@apiError OnlyOwnUser <code>403</code> code: 6007
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *academyHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)
	academy_id := ctx.UserValue("academy_id").(int)

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
	if academy_id != reqParam.Id {
		return c.Status(http.StatusForbidden).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	if err := h.academyUsecase.Patch(ctx, reqBody, reqParam.Id, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 학원 상세 조회
/**
@api {get} /academy/:id 학원 조회
@apiName getAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 상세 조회
@apiParam {Number} id 학원 id
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.name 학원 이름
@apiSuccess {String} result.callNumber 연락가능한 번호
@apiSuccess {String} result.addressRoad 도로명 주소
@apiSuccess {String} result.addressDetail 상세 주소
@apiSuccess {String} result.addressSigun 시군구
@apiSuccess {Object[]} [result.yoga] 요가
@apiSuccess {Number} result.yoga.index 요가 순서 번호
@apiSuccess {Number} result.yoga.id 요가 아이디
@apiSuccess {String} result.yoga.nameKor 요가 한글 이름
@apiSuccess {String} result.yoga.isReference 정식적으로 등록됐는지 여부
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError QueryMissing <code>400</code> code: 3001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *academyHandler) Get(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.AcademyDetailParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
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

// 학원 리스트 조회
/**
@api {get} /academy/list 학원 리스트 조회
@apiName listAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 리스트 조회
@apiQuery {Number} [pageNo=1] 페이지 번호
@apiQuery {Number} [pageSize=10] 페이지 당 문서 개수
@apiQuery {Number} [yogaIds] 필터 요가 id ,로 멀티
@apiQuery {Number} [sigunguId] 필터 시군구 id
@apiQuery {String} [academyName] 필터 학원 이름
@apiQuery {String="NAME,ID"} [orderCol="ID"] 정렬기준 컬럼 이름
@apiQuery {String="ASC,DESC"} [orderType="DESC"] 정렬기준 방법
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object[]} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.name 학원 이름
@apiSuccess {String} result.callNumber 연락가능한 번호
@apiSuccess {String} result.addressRoad 도로명 주소
@apiSuccess {String} result.addressDetail 상세 주소
@apiSuccess {String} result.addressSigun 시군구
@apiSuccess {Object[]} [result.yoga] 요가
@apiSuccess {Number} result.yoga.id 요가 아이디
@apiSuccess {String} result.yoga.nameKor 요가 한글 이름
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError QueryMissing <code>400</code> code: 3001
@apiError ValidationError <code>400</code> code: 2xxx
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *academyHandler) List(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQueries := request.NewAcademyListQueries()

	if err := fiberx.QueryParser(c, reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	if err := h.Validator.ValidateStruct(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
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
