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
	g.Delete("/:id", middleware.Auth, middleware.OnlyAcademy, handler.SoftDelete)
	g.Get("/list", handler.List)
	g.Get("/:id", handler.Get)
}

// 채용공고 생성
/**
@api {post} /recruitment 채용공고 생성
@apiName postRecruitment
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 생성
@apiHeader Authorization accessToken (Bearer)
@apiBody {Object} info
@apiBody {Boolean} info.isOpen 오픈할 건지
@apiBody {Object} insteadInfo
@apiBody {String} insteadInfo.minCareer 최소 경력
@apiBody {String} insteadInfo.pay 금액 (수업 당)
@apiBody {Object} insteadInfo.schedules
@apiBody {String} insteadInfo.schedules.startDateTime 수업 시작 일시
@apiBody {String} insteadInfo.schedules.endDateTime 수업 종료 일시
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError Conflict <code>400</code> code: 4000
@apiError ResourceUnOwned <code>400</code> code: 4010
@apiError InternalServerError <code>500</code> code: 500
*/
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

// 채용공고 수정
/**
@api {put} /recruitment/:id 채용공고 수정
@apiName putRecruitment
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 채용공고 아이디
@apiBody {Object} info
@apiBody {Boolean} info.isOpen 오픈할 건지
@apiBody {Boolean} info.isFinish 채용공고가 종료됐는지
@apiBody {Object} [insteadInfo]
@apiBody {String} insteadInfo.minCareer 최소 경력
@apiBody {String} insteadInfo.pay 금액 (수업 당)
@apiBody {Object} insteadInfo.schedules
@apiBody {String} insteadInfo.schedules.startDateTime 수업 시작 일시
@apiBody {String} insteadInfo.schedules.endDateTime 수업 종료 일시
@apiSuccess (200 or 201) {Number} code 200 or 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidationError <code>400</code> code: 2xxx
@apiError Conflict <code>400</code> code: 4000
@apiError ResourceUnOwned <code>400</code> code: 4010
@apiError InternalServerError <code>500</code> code: 500
*/
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
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
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

// 채용공고 부분 수정
/**
@api {patch} /recruitment/:id 채용공고 부분 수정
@apiName patchRecruitment
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 부분 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 채용공고 아이디
@apiBody {Object} [info]
@apiBody {Boolean} [info.isOpen] 오픈할 건지
@apiBody {Boolean} [info.isFinish] 채용공고가 종료됐는지
@apiBody {Object} [insteadInfo]
@apiBody {String} [insteadInfo.minCareer] 최소 경력
@apiBody {String} [insteadInfo.pay] 금액 (수업 당)
@apiBody {Object} [insteadInfo.schedules]
@apiBody {String} [insteadInfo.schedules.startDateTime] 수업 시작 일시
@apiBody {String} [insteadInfo.schedules.endDateTime] 수업 종료 일시
@apiSuccess (200 or 201) {Number} code 200 or 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ParamsMissing <code>400</code> code: 3002
@apiError ValidationError <code>400</code> code: 2xxx
@apiError Conflict <code>400</code> code: 4000
@apiError ResourceUnOwned <code>400</code> code: 4010
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *recruitmentHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	academyId := ctx.UserValue("academy_id").(int)

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

	isUpdated, err := h.recruitmentUsecase.Patch(ctx, reqBody, reqParam.Id, academyId)
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

// 채용공고 삭제
/**
@api {delete} /recruitment/:id 채용공고 삭제
@apiName deleteRecruitment
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 삭제
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 채용공고 아이디
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiError ParamsMissing <code>400</code> code: 3002
@apiError Conflict <code>400</code> code: 4000
@apiError ResourceUnOwned <code>400</code> code: 4010
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *recruitmentHandler) SoftDelete(c *fiber.Ctx) error {
	ctx := c.Context()
	academyId := ctx.UserValue("academy_id").(int)
	reqParam := new(request.RecruitmentDeleteParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}

	err := h.recruitmentUsecase.SoftDelete(ctx, reqParam.Id, academyId)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusOK).JSON(ex.Response{
		Code:    http.StatusOK,
		Message: "",
	})
}

// 채용공고 조회
/**
@api {get} /recruitment/:id 채용공고 조회
@apiName getRecruitment
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 조회
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 채용공고 아이디
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiSuccess {Object} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.isFinish 종료 여부
@apiSuccess {Object} result.insteadInfo 대강 정보
@apiSuccess {Number} result.insteadInfo.id 대강 아이디
@apiSuccess {String} result.insteadInfo.minCareer 최소경력
@apiSuccess {String} result.insteadInfo.pay 급여
@apiSuccess {Number} [result.insteadInfo.passerId] 합격자 아이디
@apiSuccess {Number} result.insteadInfo.applicantCount 지원자 총 명수
@apiSuccess {Object[]} [result.insteadInfo.schedules] 스케쥴
@apiSuccess {String} result.insteadInfo.schedules.startDateTime 시작 일시
@apiSuccess {String} result.insteadInfo.schedules.endDateTime 종료 일시
@apiSuccess {Object[]} [result.insteadInfo.yogas] 요가 정보
@apiSuccess {Number} result.insteadInfo.yogas.id 요가 아이디
@apiSuccess {Number} result.insteadInfo.yogas.name 요가 이름
@apiSuccess {String} result.insteadInfo.createdAt 생성일시
@apiSuccess {String} result.insteadInfo.updatedAt 업데이트일시
@apiError ParamsMissing <code>400</code> code: 3002
@apiError RecruitmentNotFound <code>400</code> code: 5006
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *recruitmentHandler) Get(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.RecruitmentGetParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}

	result, err := h.recruitmentUsecase.Get(ctx, reqParam.Id)
	if err != nil {
		return utils.NewError(c, err)
	}
	resp := response.NewRecruitmentResponse(result)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  resp,
	})
}

// 채용공고 리스트 조회
/**
@api {get} /recruitment/list 채용공고 리스트 조회
@apiName getRecruitmentList
@apiVersion 1.0.0
@apiGroup recruitment
@apiDescription 채용공고 리스트 조회
@apiHeader Authorization accessToken (Bearer)
@apiQuery {Number} [pageNo] 페이지 번호
@apiQuery {Number} [pageSize] 페이지당 문서 개수
@apiQuery {String} [startDateTime] 시작일시
@apiQuery {String} [endDateTime] 종료일시
@apiQuery {Number[]} [yogaIds] 요가 아이디
@apiQuery {Number[]} [sigunguIds] 시군구 아이디
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiSuccess {Object[]} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.academyName 학원 이름
@apiSuccess {String[]} result.yogas 요가 이름
@apiSuccess {String} result.sigungu 시군구 이름
@apiSuccess {String[]} result.startDateTimes 시작 일시
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError ParamsMissing <code>400</code> code: 3002
@apiError RecruitmentNotFound <code>400</code> code: 5006
@apiError InternalServerError <code>500</code> code: 500
*/
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
