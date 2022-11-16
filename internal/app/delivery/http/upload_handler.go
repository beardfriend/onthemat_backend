package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport"
	"onthemat/internal/app/usecase"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type uploadHandler struct {
	uploadUseCase usecase.UploadUsecase
	Validator     validatorx.Validator
}

func NewUploadHandler(
	middleware *middlewares.MiddleWare,
	uploadUseCase usecase.UploadUsecase,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &uploadHandler{
		Validator:     validator,
		uploadUseCase: uploadUseCase,
	}

	g := router.Group("/upload")
	// 이미지 업로드
	g.Post("/:purpose", middleware.Auth, handler.Upload)
}

// 이미지 업로드
/**
@api {get} /upload/:purpose 이미지 업로드
@apiName acessTokenRefresh
@apiVersion 1.0.0
@apiGroup upload
@apiDescription 이미지를 업로드하는 API

@apiHeader {String} Authorization Bearer 엑세스토큰
@apiParam {String="profile,logo"} purpose 업로드 이후 사용할 목적

@apiSuccessExample Success-Response:
HTTP/1.1 201 Created
{
	"code": 201,
	"message": ""
}
@apiErrorExample Error-Response:
HTTP/1.1 400 Bad Request
{
	"code": 400,
	"message": "bad request",
	"detail": "파라메터를 확인해주세요."
}

HTTP/1.1 400 Bad Request
{
    "code": 400,
    "message": "bad request",
    "details": [
        {
            "purpose": "oneof"
        }
    ]
}

HTTP/1.1 400 Bad Request
{
    "code": 400,
    "message": "bad request",
    "details": "form key를 확인해주세요."
}

HTTP/1.1 400 Bad Request
{
    "code": 400,
    "message": "bad request",
    "details": "form key를 확인해주세요."
}

HTTP/1.1 400 Bad Request
{
    "code": 400,
    "message": "bad request",
    "details": "사용할 수 없는 확장자입니다."
}

HTTP/1.1 500 Internal Server Error
{
	"code": 500,
	"message": "internal server error",
	"detail": "일시적인 에러가 발생했습니다."
}
*/
func (h *uploadHandler) Upload(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqParams := new(transport.UploadParams)
	if err := c.ParamsParser(reqParams); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrParamsMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqParams); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(ex.NewInvalidInputError(err))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(ex.NewHttpError(ex.ErrFormDataKeyUnavailable, "key name is file"))
	}

	if err := h.uploadUseCase.Upload(ctx, file, reqParams, userId); err != nil {
		code, json := ex.ParseHttpError(err)
		return c.Status(code).JSON(json)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}
