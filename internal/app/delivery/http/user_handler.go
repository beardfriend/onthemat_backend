package http

import (
	_ "database/sql"
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport/response"
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
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} [result.email] 이메일
@apiSuccess {String} [result.nickname] 닉네임
@apiSuccess {String="kakao,google,naver"} [result.social_name] 소셜로그인 타입
@apiSuccess {String} [result.social_key] 소셜 고유 Key
@apiSuccess {String="academy,teacher,superAdmin"} [result.type] 유저 타입
@apiSuccess {String} [result.phone_num] 이메일
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError UserNotFound <code>404</code> code: 5001
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := c.Context().UserValue("user_id").(int)

	u, err := h.UserUseCase.GetMe(ctx, userId)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewUserMeResponse(u)
	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    200,
		Message: "",
		Result:  resp,
	})
}
