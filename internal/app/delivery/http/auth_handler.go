package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/transport/response"
	"onthemat/internal/app/utils"

	"onthemat/internal/app/usecase"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	AuthUseCase usecase.AuthUseCase
	Validator   validatorx.Validator
}

func NewAuthHandler(
	authUseCase usecase.AuthUseCase,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &authHandler{
		AuthUseCase: authUseCase,
		Validator:   validator,
	}
	g := router.Group("/auth")
	// 소셜로그인 리디렉션
	g.Get("/:socialName/url", handler.SocialUrl)
	// 소셜로로그인 콜백
	g.Get("/:socialName/callback", handler.SocialCallback)
	// 회원가입
	g.Post("/signup", handler.SignUp)
	// 로그인
	g.Post("/login", handler.Login)
	// 소셜 회원가입
	g.Patch("/social/signup", handler.SocialSignUp)
	// 임시 비밀번호 발급
	g.Get("/temp-password", handler.SendTempPassword)
	// 이메일 중복체크
	g.Get("/check-email", handler.CheckDuplicatedEmail)
	// 이메일 인증
	g.Get("/verify-email", handler.VerifiyEmail)
	// Access Token 리프레쉬
	g.Get("/token/refresh", handler.Refresh)
}

// 소셜로그인 리디렉션
/**
@api {get} /auth/:socialName/url 소셜로그인 URL
@apiName socialLoginRedirection
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 소셜로그인 URL
@apiParam {String="naver,kakao,google"} socialName 소셜 로그인 타입
@apiSuccessExample Success-Response:
HTTP/1.1 302 OK
*/
func (h *authHandler) SocialUrl(c *fiber.Ctx) error {
	ctx := c.Context()
	reqParam := new(request.SocialUrlParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqParam); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	url, err := h.AuthUseCase.SocialLoginRedirectUrl(ctx, reqParam.SocialName)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Redirect(url)
}

// 소셜로그인 콜백
/**
@api {get} /auth/:socialName/callback 소셜 로그인 콜백
@apiName Socialcallback
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 소셜 Callback
@apiParam {String="naver,kakao,google"} socialName 소셜 로그인 타입
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {String} result.accessToken 엑세스 토큰
@apiSuccess {String} result.accessTokenexpiredAt 엑세스 토큰 만료일시
@apiSuccess {String} result.refreshToken 리프레쉬 토큰
@apiSuccess {String} result.refreshTokenExpiredAt 리프레쉬 토큰 만료일시
@apiError QueryStringMissing <code>400</code> code: 3001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) SocialCallback(c *fiber.Ctx) error {
	ctx := c.Context()

	reqParam := new(request.SocialCallbackParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqParam); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	data, err := h.AuthUseCase.SocialLogin(ctx, reqParam.SocialName, code)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(200).
		JSON(ex.ResponseWithData{
			Code:    200,
			Message: "",
			Result:  data,
		})
}

// 회원가입
/**
@api {post} /auth/signup 일반 회원가입
@apiName signup
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 일반 회원가입 API

@apiBody {String} email 이메일
@apiBody {String} password 비밀번호
@apiBody {String} nickname 닉네임
@apiBody {Boolean} termAgree 약관동의여부

@apiSuccessExample Success-Response:
HTTP/1.1 201 OK
{
	"code": 201,
	"message": ""
}

@apiErrorExample Error-Response:
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
            "nickname": "max"
        },
        {
            "termAgree": "required"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2001,
    "message": "유효하지 않은 패스워드입니다.",
    "details": [
        {
            "password": "max"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}
HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.SignUpBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.Validator.ValidateStruct(body); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.SignUp(ctx, body); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(201).
		JSON(ex.Response{
			Code:    201,
			Message: "",
		})
}

// 소셜 회원가입
/**
@api {patch} /auth/social/signup 소셜 회원가입
@apiName socialSingup
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 소셜회원가입 API

@apiBody {Number} userId 유저의 Primary Key
@apiBody {String} email 이메일


@apiSuccessExample Success-Response:
HTTP/1.1 201 OK
{
	"code": 201,
	"message": ""
}

@apiErrorExample Error-Response:
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
            "userId": "required"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}

HTTP/1.1 404 Not Found
{
    "code": 5001,
    "message": "존재하지 않는 유저입니다.",
    "details": null
}

HTTP/1.1 409 Conflict
{
    "code": 4001,
    "message": "이미 존재하는 이메일입니다.",
    "details": null
}

HTTP/1.1 409 Conflict
{
    "code": 4004,
    "message": "이미 회원의 이메일이 등록됐습니다.",
    "details": null
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) SocialSignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.SocialSignUpBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(body); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.SocialSignUp(ctx, body); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(201).
		JSON(ex.Response{
			Code:    201,
			Message: "",
		})
}

// 로그인
/**
@api {post} /auth/login 일반 로그인
@apiName login
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 일반 로그인 API

@apiBody {String} email 이메일
@apiBody {String} password 비밀번호

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": ""
}

@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
    "code": 3000,
    "message": "JSON을 입력해주세요.",
    "details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2001,
    "message": "유효하지 않은 패스워드입니다.",
    "details": [
        {
            "password": "max"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}

HTTP/1.1 404 Not Found
{
    "code": 5001,
    "message": "존재하지 않는 유저입니다.",
    "details": "이메일 혹은 비밀번호를 다시 확인해주세요."
}

HTTP/1.1 401 Unauthorized
{
    "code": 6001,
    "message": "유저의 이메일 인증이 필요합니다.",
    "details": null
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(transport.LoginBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	if err := h.Validator.ValidateStruct(body); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	data, err := h.AuthUseCase.Login(ctx, body)
	if err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(200).
		JSON(ex.ResponseWithData{
			Code:    200,
			Message: "",
			Result:  data,
		})
}

// 이메일 중복체크
/**
@api {get} /auth/check-email 이메일 중복체크
@apiName duplicatedEmail
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 이메일 중복체크 API

@apiQuery {String} email 유저의 email

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": ""
}
@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 3001,
	"message": "쿼리스트링을 입력해주세요.",
	"details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}

HTTP/1.1 409 Conflict
{
	"code": 4001,
	"message": "이미 존재하는 이메일입니다.",
	"details": null
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) CheckDuplicatedEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(transport.CheckDuplicatedEmailQueries)

	if err := c.QueryParser(queries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(queries); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.CheckDuplicatedEmail(ctx, queries.Email); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(200).
		JSON(ex.Response{
			Code:    200,
			Message: "",
		})
}

// 임시 비밀번호 발급
/**
@api {get} /auth/temp-password 임시 비밀번호 이메일 발송
@apiName tempPassword
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 유저의 이메일 계정으로 임시비밀번호 발송하는 API

@apiQuery {String} email 유저의 email

@apiSuccessExample Success-Response:
HTTP/1.1 202 Accepted
{
	"code": 202,
	"message": ""
}

@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 3001,
	"message": "쿼리스트링을 입력해주세요.",
	"details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}

HTTP/1.1 404 Not Found
{
    "code": 5001,
    "message": "존재하지 않는 유저입니다.",
    "details": "존재하지 않는 이메일입니다."
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) SendTempPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(transport.SendTempPasswordQueries)

	if err := c.QueryParser(queries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(queries); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.SendEmailResetPassword(ctx, queries.Email); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusAccepted).
		JSON(ex.Response{
			Code:    200,
			Message: "",
		})
}

// 이메일 인증
/**
@api {get} /auth/verify-email 이메일 인증
@apiName verifyEmail
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 이메일 인증하는 API

@apiQuery {String} email 유저의 email
@apiQuery {String} key 이메일에 포함된 key값

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": ""
}
@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 3001,
	"message": "쿼리스트링을 입력해주세요.",
	"details": null
}

HTTP/1.1 400 Bad Request
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
	"code": 3006,
	"message": "사용할 수 없는 인증 키입니다.",
	"details": null
}

HTTP/1.1 401 Authentication UnAuthorization
{
	"code": 6003,
	"message": "이메일 인증 시간이 만료되었습니다.",
	"details": null
}

HTTP/1.1 404 Not Found
{
    "code": 5001,
    "message": "존재하지 않는 유저입니다.",
    "details": null
}

HTTP/1.1 409 Conflict
{
	"code": 4002,
	"message": "이미 이메일 인증이 완료됐습니다.",
	"details": null
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}
*/
func (h *authHandler) VerifiyEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(transport.VerifyEmailQueries)

	if err := c.QueryParser(queries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrQueryStringMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(queries); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	if err := h.AuthUseCase.VerifiyEmail(ctx, queries.Email, queries.Key, queries.IssuedAt); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(200).
		JSON(ex.Response{
			Code:    200,
			Message: "",
		})
}

// Access Token 리프레쉬
/**
@api {get} /auth/token/refresh Access 토큰 리프레쉬
@apiName acessTokenRefresh
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 엑세스 토큰을 재발급하는 API

@apiHeader {String} Authorization Bearer 리프레쉬토큰

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
    "code": 200,
    "message": "",
    "result": {
        "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiZGE1MmE3NTgtNmZjZS00MDhkLTlmYWQtN2E0YTk4Njc1YmQ3IiwiVXNlcklkIjoxLCJMb2dpblR5cGUiOiJub3JtYWwiLCJVc2VyVHlwZSI6IiIsImlzcyI6Im9uZVRoZU1hdCIsImV4cCI6MTcyODkwMzkxMCwiaWF0IjoxNjY4OTAzOTEwfQ.945I_QAC2uK7f4YQdWQKn_z0RotL3FlLs9J6XbztBmg",
        "accessTokenExpiredAt": "2024-10-14T20:05:10"
    }
}
@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 3005,
	"message": "Authorization 헤더 형식을 확인해주세요.",
	"details": null
}

HTTP/1.1 400 Bad Request
{
	"code": 3007,
	"message": "유효하지 않은 토큰입니다.",
	"details": null
}

HTTP/1.1 401 Authentication UnAuthorization
{
	"code": 6002,
	"message": "토큰이 만료되었습니다.",
	"details": null
}

HTTP/1.1 404 Not Found
{
	"code": 5001,
	"message": "존재하지 않는 유저입니다.",
	"details": null
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "일시적인 에러가 발생했습니다.",
	"details": null
}

*/
func (h *authHandler) Refresh(c *fiber.Ctx) error {
	ctx := c.Context()

	authorizationHeader := c.Request().Header.Peek("Authorization")

	data, err := h.AuthUseCase.Refresh(ctx, authorizationHeader)
	if err != nil {
		return utils.NewError(c, err)
	}

	res := response.NewRefreshResponse(data)

	return c.Status(200).
		JSON(ex.ResponseWithData{
			Code:    200,
			Message: "",
			Result:  res,
		})
}
