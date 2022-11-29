package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
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
	reqParam := new(request.AuthSocialUrlParam)

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

	reqParam := new(request.AuthSocialCallbackParam)
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
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError ErrJsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError PasswordInvalid <code>400</code> code: 2001
@apiError EmailInvalid <code>400</code> code: 2002
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) SignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(request.AuthSignUpBody)
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
@apiBody {Number} userId 유저 아이디
@apiBody {String} email 이메일
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError ErrJsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError PasswordInvalid <code>400</code> code: 2001
@apiError UserNotFound <code>404</code> code: 5001
@apiError UserEmailAlreadyExist <code>409</code> code: 4001
@apiError UserEmailAlreadyRegisted <code>409</code> code: 4004
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) SocialSignUp(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(request.AuthSocialSignUpBody)
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
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {String} result.accessToken 엑세스 토큰
@apiSuccess {String} result.accessTokenexpiredAt 엑세스 토큰 만료일시
@apiSuccess {String} result.refreshToken 리프레쉬 토큰
@apiSuccess {String} result.refreshTokenExpiredAt 리프레쉬 토큰 만료일시
@apiError PasswordInvalid <code>400</code> code: 2001
@apiError EmailInvalid <code>400</code> code: 2002
@apiError UserEmailUnauthorization <code>401</code> code: 6001
@apiError UserNotFound <code>404</code> code: 5001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	body := new(request.AuthLoginBody)
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
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiError QueryStringMissing <code>400</code> code: 3001
@apiError EmailInvalid <code>400</code> code: 2002
@apiError UserEmailAlreadyExist <code>409</code> code: 4001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) CheckDuplicatedEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(request.AuthCheckDuplicatedEmailQuery)

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
@apiSuccess (202) {Number} code 202
@apiSuccess (202) {String} message ""
@apiError QueryStringMissing <code>400</code> code: 3001
@apiError EmailInvalid <code>400</code> code: 2002
@apiError UserNotFound <code>404</code> code: 5001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) SendTempPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(request.AuthTempPasswordQuery)

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
			Code:    202,
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
@apiSuccess (200) {Number} code 200
@apiSuccess (200) {String} message ""
@apiError QueryStringMissing <code>400</code> code: 3001
@apiError EmailInvalid <code>400</code> code: 2002
@apiError RandomKeyForEmailVerfiyUnavailable <code>400</code> code: 3006
@apiError EmailForVerifyExpired <code>401</code> code: 6003
@apiError UserNotFound <code>404</code> code: 5001
@apiError UserEmailAlreadyVerfied <code>409</code> code: 4002
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *authHandler) VerifiyEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	queries := new(request.AuthVerifyEmailQueries)

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
@apiHeader {String} Authorization 리프레쉬토큰(Bearer)
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {String} result.accessToken 엑세스 토큰
@apiSuccess {String} result.accessTokenexpiredAt 엑세스 토큰 만료일시
@apiError AuthorizationHeaderFormatUnavailable <code>400</code> code: 3005
@apiError TokenInvalid <code>400</code> code: 3007
@apiError TokenExpired <code>401</code> code: 6002
@apiError UserNotFound <code>404</code> code: 5001
@apiError InternalServerError <code>500</code> code: 500
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
