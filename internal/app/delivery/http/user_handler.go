package http

import (
	_ "database/sql"
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
	Router      fiber.Router
	Middleware  middlewares.MiddleWare
}

func NewUserHandler(
	middleware middlewares.MiddleWare,
	userUseCase usecase.UserUseCase,
	router fiber.Router,
) {
	handler := &UserHandler{
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
@apiName getUser
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
        "email": "beardfriend21@naver.com",
        "nickname": "nick",
        "social_name": null,
        "social_key": null,
        "type": null,
        "phone_num": null,
        "created_at": "2022-11-20T09:06:20",
        "last_login_at": "2022-11-20T09:06:20"
    }
}

HTTP/1.1 500 Internal Server Error
{
    "code": 500,
    "message": "일시적인 에러가 발생했습니다.",
    "details": null
}
*/

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := c.Context().UserValue("user_id").(int)
	u, err := h.UserUseCase.GetMe(ctx, userId)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := transport.NewUserMeResponse(u)
	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    200,
		Message: "",
		Result:  resp,
	})
}
