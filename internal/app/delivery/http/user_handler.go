package http

import (
	"onthemat/internal/app/common"
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	UserUseCase usecase.UserUseCase
	Router      fiber.Router
	Middleware  *middlewares.MiddleWare
}

func NewUserHandler(
	middleware *middlewares.MiddleWare,
	userUseCase usecase.UserUseCase,
	router fiber.Router,
) {
	handler := &userHandler{
		UserUseCase: userUseCase,
		Middleware:  middleware,
	}
	g := router.Group("/user")
	/**
	@api {get} /user/me 내 정보 조회
	@apiName naverCallback
	@apiVersion 1.0.0
	@apiGroup user
	@apiDescription 내 정보를 조회하는 API
	@apiHeader Authorization accessToken (Bearer)

	@apiSuccessExample Success-Response:
	HTTP/1.1 200 OK
	{
		"code": 200,
		"message": "",
		"result": {
			"id": 1,
			"email": "beardfriend21@gmail.com",
			"nickname": "",
			"created_at": "2022-11-07T10:20:21.797615+09:00",
			"social_name": "",
			"social_key": "",
			"type": "",
			"phone_num": "",
			"last_login_at": "2022-11-07T10:19:53.805358+09:00"
		}
	}

	HTTP/1.1 500 Internal Server Error
	{
		"code": 500,
		"message": "internal server error",
		"detail": "일시적인 에러가 발생했습니다."
	}
	*/
	g.Get("/me", middleware.Auth, handler.GetMe)
}

func (h *userHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()

	u, err := h.UserUseCase.GetMe(ctx, c.Context().UserValue("user_id").(int))
	if err != nil {
		code, json := ex.ParseHttpError(err)
		return c.Status(code).JSON(json)
	}

	resp := transport.NewUserMeResponse(u)
	return c.JSON(common.ResponseWithData{
		Code:    200,
		Message: "",
		Result:  resp,
	})
}
