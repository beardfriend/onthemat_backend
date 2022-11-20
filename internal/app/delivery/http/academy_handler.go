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
	Validator      validatorx.Validator
}

func NewAcademyHandler(
	middleware *middlewares.MiddleWare,
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
	// 학원 정보 수정
	g.Put("", middleware.Auth, middleware.OnlyAcademy, handler.Update)
	// 학원 리스트
	g.Get("/list", handler.List)
	// 학원 상세조회
	g.Get("/:id", handler.Detail)
}

// 학원 등록
/**
@api {post} /academy 학원 등록
@apiName postAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 회원가입
@apiHeader Authorization accessToken (Bearer)

@apiSuccessExample Success-Response:
HTTP/1.1 201 Created
{
    "code": 200,
    "message": ""
}

HTTP/1.1 400 Bad Request
{
    "code": 3000,
    "message": "JSON을 입력해주세요.",
    "details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2000,
    "message": "유효하지 않은 요청값들이 존재합니다.",
    "details": [
        {
            "addressDong": "required"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2004,
    "message": "유효하지 않은 휴대폰번호입니다.",
    "details": [
        {
            "callNumber": "phoneNumNoDash"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2003,
    "message": "유효하지 않은 url입니다.",
    "details": [
        {
            "logoUrl": "urlStartHttpHttps"
        }
    ]
}

HTTP/1.1 409 Conflict
{
    "code": 4003,
    "message": "이미 회원 유형이 등록됐습니다.",
    "details": null
}


HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *academyHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(transport.AcademyCreateRequestBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
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
@api {put} /academy 학원 정보 수정
@apiName putAcademy
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 정보 수정
@apiHeader Authorization accessToken (Bearer)

@apiSuccessExample Success-Response:
HTTP/1.1 201 Created
{
    "code": 200,
    "message": ""
}

HTTP/1.1 400 Bad Request
{
    "code": 3000,
    "message": "JSON을 입력해주세요.",
    "details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2000,
    "message": "유효하지 않은 요청값들이 존재합니다.",
    "details": [
        {
            "addressDong": "required"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2004,
    "message": "유효하지 않은 휴대폰번호입니다.",
    "details": [
        {
            "callNumber": "phoneNumNoDash"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2003,
    "message": "유효하지 않은 url입니다.",
    "details": [
        {
            "logoUrl": "urlStartHttpHttps"
        }
    ]
}

HTTP/1.1 404 Not Found
{
    "code": 5003,
    "message": "존재하지 않는 학원입니다.",
    "details": null
}


HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *academyHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(transport.AcademyUpdateRequestBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.academyUsecase.Update(ctx, reqBody, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 학원 상세 조회
/**
@api {get} /academy/:id 학원 상세 조회
@apiName getAcademyDetail
@apiVersion 1.0.0
@apiGroup academy
@apiDescription 학원 상세보기
@apiHeader Authorization accessToken (Bearer)

@apiParam {Number} id 학원 아이디

@apiSuccessExample Success-Response:
HTTP/1.1 200
{
    "code": 200,
    "message": "",
    "result": {
        "id": 1,
        "name": "학원이름이모르지",
        "callNumber": "01043226632",
        "addressRoad": "서울시 양천구 도곡로 25길 10-2",
        "addressDetail": "",
        "addressSigun": "서울시",
        "addressGu": "양천구",
        "addressX": "13230.123",
        "addressY": "123.232",
        "createdAt": "2022-11-20T09:26:19",
        "updatedAt": "2022-11-20T09:31:20"
    }
}

HTTP/1.1 400 Bad Request
{
    "code": 3002,
    "message": "파라메터를 입력해주세요.",
    "details": null
}

HTTP/1.1 404 Not Found
{
    "code": 5003,
    "message": "존재하지 않는 학원입니다.",
    "details": null
}


HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *academyHandler) Detail(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(transport.AcademyDetailParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}
	if err := h.Validator.ValidateStruct(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	academy, err := h.academyUsecase.Get(ctx, reqParam.Id)
	if err != nil {
		return utils.NewError(c, err)
	}

	response := transport.NewAcademyDetailResponse(academy)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  response,
	})
}

func (h *academyHandler) List(c *fiber.Ctx) error {
	ctx := c.Context()

	reqQueries := transport.NewAcademyListQueries()

	if err := c.QueryParser(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	if err := h.Validator.ValidateStruct(reqQueries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	academies, pagination, err := h.academyUsecase.List(ctx, reqQueries)
	if err != nil {
		return utils.NewError(c, err)
	}

	response := transport.NewAcademyListResponse(academies)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithPagination{
		Code:       http.StatusOK,
		Message:    "",
		Result:     response,
		Pagination: pagination,
	})
}
