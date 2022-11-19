package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/model"
	"onthemat/internal/app/transport"
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
	// 카카오 리디렉션
	g.Get("/kakao", handler.Kakao)
	// 카카오 콜백
	g.Get("/kakao/callback", handler.KakaoCallBackToken)
	// 구글 리디렉션
	g.Get("/google", handler.Google)
	// 구글 콜백
	g.Get("/google/callback", handler.GoogleCallBackToken)
	// 네이버 리디렉션
	g.Get("/naver", handler.Naver)
	// 네이버 콜백
	g.Get("/naver/callback", handler.NaverCallBackToken)
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

// 카카오 리디렉션
/**
@api {get} /auth/kakao 카카오 리디렉션
@apiName kakao
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 카카오 리디렉션 URL

@apiSuccessExample Success-Response:
HTTP/1.1 302 OK
*/
func (h *authHandler) Kakao(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.KakaoRedirectUrl(ctx))
}

// 카카오 콜백
/**
@api {get} /auth/kakao/callback 카카오 로그인 콜백 URL
@apiName kakaoCallback
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 카카오 Callback

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": "",
	"result": {
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY3ODAzMTAyLCJpYXQiOjE2Njc4MDIyMDJ9.wFaNMotM7E38mM_Rcyk5GlAe7WTUX-zJv9CPGgixpds",
		"accessTokenexpiredAt": "2022-11-07T15:38:22.270238+09:00",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY5MDExODAyLCJpYXQiOjE2Njc4MDIyMDJ9.mXJ4QM19pHrM_4pNFVF1d1PnCYBLTRR4EaYc70O2N88",
		"refreshTokenExpiredAt": "2022-11-21T15:23:22.270239+09:00"
	}
}

@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
    "code": 3001,
    "message": "쿼리스트링을 입력해주세요.",
    "details": null
}

HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *authHandler) KakaoCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	data, err := h.AuthUseCase.SocialLogin(ctx, model.KakaoSocialType, code)
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

// 구글 리디렉션
/**
@api {get} /auth/google 구글 리디렉션
@apiName google
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 구글 리디렉션 URL

@apiSuccessExample Success-Response:
HTTP/1.1 302 OK
*/
func (h *authHandler) Google(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.GoogleRedirectUrl(ctx))
}

// 구글 콜백
/**
@api {get} /auth/google/callback 구글 로그인 콜백 URL
@apiName googleCallback
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 구글 Callback

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": "",
	"result": {
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY3ODAzMTAyLCJpYXQiOjE2Njc4MDIyMDJ9.wFaNMotM7E38mM_Rcyk5GlAe7WTUX-zJv9CPGgixpds",
		"accessTokenexpiredAt": "2022-11-07T15:38:22.270238+09:00",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY5MDExODAyLCJpYXQiOjE2Njc4MDIyMDJ9.mXJ4QM19pHrM_4pNFVF1d1PnCYBLTRR4EaYc70O2N88",
		"refreshTokenExpiredAt": "2022-11-21T15:23:22.270239+09:00"
	}
}

@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
    "code": 3001,
    "message": "쿼리스트링을 입력해주세요.",
    "details": null
}

HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *authHandler) GoogleCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(ex.NewBadRequestError(ex.ErrQueryStringMissing, nil))
	}

	data, err := h.AuthUseCase.SocialLogin(ctx, model.GoogleSocialType, code)
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

// 네이버 리디렉션
/**
@api {get} /auth/naver 네이버 리디렉션
@apiName naver
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 네이버 리디렉션 URL

@apiSuccessExample Success-Response:
HTTP/1.1 302 OK
*/
func (h *authHandler) Naver(c *fiber.Ctx) error {
	ctx := c.Context()

	return c.Redirect(h.AuthUseCase.NaverRedirectUrl(ctx))
}

// 네이버 콜백
/**
@api {get} /auth/naver/callback 네이버 로그인 콜백 URL
@apiName naverCallback
@apiVersion 1.0.0
@apiGroup auth
@apiDescription 네이버 Callback

@apiSuccessExample Success-Response:
HTTP/1.1 200 OK
{
	"code": 200,
	"message": "",
	"result": {
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY3ODAzMTAyLCJpYXQiOjE2Njc4MDIyMDJ9.wFaNMotM7E38mM_Rcyk5GlAe7WTUX-zJv9CPGgixpds",
		"accessTokenexpiredAt": "2022-11-07T15:38:22.270238+09:00",
		"refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY5MDExODAyLCJpYXQiOjE2Njc4MDIyMDJ9.mXJ4QM19pHrM_4pNFVF1d1PnCYBLTRR4EaYc70O2N88",
		"refreshTokenExpiredAt": "2022-11-21T15:23:22.270239+09:00"
	}
}

@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
    "code": 3001,
    "message": "쿼리스트링을 입력해주세요.",
    "details": null
}

HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *authHandler) NaverCallBackToken(c *fiber.Ctx) error {
	ctx := c.Context()

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrQueryStringMissing, nil))
	}

	data, err := h.AuthUseCase.SocialLogin(ctx, model.NaverSocialType, code)
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
{
    "code": 2001,
    "message": "유효하지 않은 패스워드입니다.",
    "details": [
        {
            "password": "max"
        }
    ]
}
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
{
    "code": 2000,
    "message": "유효하지 않은 요청값들이 존재합니다.",
    "details": [
        {
            "userId": "required"
        }
    ]
}
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
{
    "code": 2001,
    "message": "유효하지 않은 패스워드입니다.",
    "details": [
        {
            "password": "max"
        }
    ]
}
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
{
    "code": 2002,
    "message": "유효하지 않은 이메일입니다.",
    "details": [
        {
            "email": "email"
        }
    ]
}
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
		"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVdWlkIjoiNjE5YWUxYTYtN2YyNy00NDZmLTkzZGUtNDBjNjJkM2MwOWU3IiwiVXNlcklkIjowLCJMb2dpblR5cGUiOiJrYWthbyIsIlVzZXJUeXBlIjoiIiwiaXNzIjoib25lVGhlTWF0IiwiZXhwIjoxNjY3ODAzMTAyLCJpYXQiOjE2Njc4MDIyMDJ9.wFaNMotM7E38mM_Rcyk5GlAe7WTUX-zJv9CPGgixpds",
		"accessTokenexpiredAt": "2022-11-07T15:38:22.270238+09:00"
	}
}
@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 3005,
	"message": "Authorization 헤더 형식을 확인해주세요.",
	"details": null
}
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

	return c.Status(200).
		JSON(ex.ResponseWithData{
			Code:    200,
			Message: "",
			Result:  data,
		})
}
