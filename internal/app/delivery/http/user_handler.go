package http

import (
	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"

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
	// 유저 정보 조회
	g.Get("/me", middleware.Auth, handler.GetMe)
}

// 유저 정보 조회
/**
@api {get} /user/me 유저 정보 조회
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
        "id": 2,
        "email": "beardfriend22@gmail.com",
        "nickname": "nick",
        "social_name": null,
        "social_key": null,
        "type": null,
        "phone_num": null,
        "created_at": "2022-11-11T17:21:41+09:00",
        "last_login_at": "2022-11-11T17:20:44+09:00"
    }
}

HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/
func (h *userHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()

	u, err := h.UserUseCase.GetMe(ctx, c.Context().UserValue("user_id").(int))
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := transport.NewUserMeResponse(u)
	return c.JSON(ex.ResponseWithData{
		Code:    200,
		Message: "",
		Result:  resp,
	})
}
