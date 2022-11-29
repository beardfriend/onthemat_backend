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

type uploadHandler struct {
	uploadUseCase usecase.UploadUsecase
	Validator     validatorx.Validator
	middleware    middlewares.MiddleWare
}

func NewUploadHandler(
	middleware middlewares.MiddleWare,
	uploadUseCase usecase.UploadUsecase,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &uploadHandler{
		Validator:     validator,
		uploadUseCase: uploadUseCase,
		middleware:    middleware,
	}

	g := router.Group("/upload")
	// 이미지 업로드
	g.Post("/:purpose", middleware.Auth, handler.Upload)
}

// 이미지 업로드
/**
@api {get} /upload/:purpose 이미지 업로드
@apiName uploadImage
@apiVersion 1.0.0
@apiGroup upload
@apiDescription 이미지를 업로드하는 API
@apiHeader {String} Authorization accessToken (Bearer)
@apiParam {String="profile,logo"} purpose 업로드 이후 사용할 목적
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError ValidationError <code>400</code> code: 2xxx
@apiError ImageExtensionUnavailable <code>400</code> code: 3003
@apiError FormDataKeyUnavailable <code>400</code> code: 3004
@apiError InternalServerError <code>500</code> code: 500
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
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}
