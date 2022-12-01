package http

import (
	"net/http"

	ex "onthemat/internal/app/common"
	"onthemat/internal/app/delivery/middlewares"
	"onthemat/internal/app/transport/request"
	"onthemat/internal/app/transport/response"
	"onthemat/internal/app/usecase"
	"onthemat/internal/app/utils"
	"onthemat/pkg/fiberx"
	"onthemat/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type teacherHandler struct {
	teacherUsecase usecase.TeacherUsecase
	Validator      validatorx.Validator
}

func NewTeacherHandler(
	middleware middlewares.MiddleWare,
	teacherUsecase usecase.TeacherUsecase,
	validator validatorx.Validator,
	router fiber.Router,
) {
	handler := &teacherHandler{
		teacherUsecase: teacherUsecase,
		Validator:      validator,
	}

	g := router.Group("/teacher")
	// 선생님 생성
	g.Post("", middleware.Auth, handler.Create)
	// 선생님 수정
	g.Put("/:id", middleware.Auth, middleware.OnlyTeacher, handler.Update)
	// 선생님 부분 수정
	g.Patch("/:id", middleware.Auth, middleware.OnlyTeacher, handler.Patch)
	// 선생님 상세 조회
	g.Get("/:id", middleware.Auth, handler.Get)
}

// 선생님 생성
/**
@api {post} /teacher 선생님 생성
@apiName postTeacher
@apiVersion 1.0.0
@apiGroup teacher
@apiDescription 선생님 생성
@apiHeader Authorization accessToken (Bearer)
@apiBody {Object} info
@apiBody {String} info.name 선생님 이름
@apiBody {String} [info.profileImageUrl] 대표 이미지 주소
@apiBody {Number} [info.age] 나이
@apiBody {String} [info.introduce] 간단한 자기소개
@apiBody {Boolean} info.isProfileOpen 프로필 공개여부
@apiBody {Object[]} [workExperiences] 근무 경험
@apiBody {String} workExperiences.academyName 근무지 학원 이름
@apiBody {String} workExperiences.workStartAt 근무시작 일시
@apiBody {String} [workExperiences.WorkEndAt] 근무종료 일시
@apiBody {String} [workExperiences.description] 근무에 대한 설명
@apiBody {Object[]} [certifications] 자격증
@apiBody {String} certifications.agencyName 취득 기관 이름
@apiBody {String} certifications.imageUrl 자격증 이미지 주소
@apiBody {String} certifications.classStartAt 수업시작 일시
@apiBody {String} [certifications.classEndAt] 수업종료 일시
@apiBody {String} [certifications.description] 자격증에 대한 설명
@apiBody {Object[]} [yogaRaws] 직접 등록한 요가 (검색 시스템에 나오지 않는 경우 회원이 취급하는 요가 종류 등록)
@apiBody {String} yogaRaws.name 요가 이름
@apiBody {int[]} [yogaIds] 수업 가능한 요가 아이디
@apiBody {int[]} [sigunguIds] 활동가능한 지역 아이디
@apiSuccess (201) {Number} code 201
@apiSuccess (201) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError UserNotFound <code>404</code> code: 5001
@apiError UserTypeAlreadyRegisted <code>409</code> code: 4003
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *teacherHandler) Create(c *fiber.Ctx) error {
	ctx := c.Context()

	userId := ctx.UserValue("user_id").(int)

	reqBody := new(request.TeacherCreateBody)

	if err := fiberx.BodyParser(c, reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, err.Error()))
	}

	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	if err := h.teacherUsecase.Create(ctx, reqBody, userId); err != nil {
		return utils.NewError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(ex.Response{
		Code:    http.StatusCreated,
		Message: "",
	})
}

// 선생님 수정
/**
@api {put} /teacher/:id 선생님 수정
@apiName putTeacher
@apiVersion 1.0.0
@apiGroup teacher
@apiDescription 선생님 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 선생님 아이디
@apiBody {Object} info
@apiBody {String} info.name 선생님 이름
@apiBody {String} [info.profileImageUrl] 대표 이미지 주소
@apiBody {Number} [info.age] 나이
@apiBody {String} [info.introduce] 간단한 자기소개
@apiBody {Boolean} info.isProfileOpen 프로필 공개여부
@apiBody {Object[]} [workExperiences] 근무 경험
@apiBody {Number} workExperiences.id 아이디
@apiBody {String} workExperiences.academyName 근무지 학원 이름
@apiBody {String} workExperiences.workStartAt 근무시작 일시
@apiBody {String} [workExperiences.WorkEndAt] 근무종료 일시
@apiBody {String} [workExperiences.description] 근무에 대한 설명
@apiBody {Object[]} [certifications] 자격증
@apiBody {Number} certifications.id 아이디
@apiBody {String} certifications.agencyName 취득 기관 이름
@apiBody {String} certifications.imageUrl 자격증 이미지 주소
@apiBody {String} certifications.classStartAt 수업시작 일시
@apiBody {String} [certifications.classEndAt] 수업종료 일시
@apiBody {String} [certifications.description] 자격증에 대한 설명
@apiBody {Object[]} [yogaRaws] 직접 등록한 요가 (검색 시스템에 나오지 않는 경우 회원이 취급하는 요가 종류 등록)
@apiBody {Number} yogaRaws.id 아이디
@apiBody {String} yogaRaws.name 요가 이름
@apiBody {int[]} [yogaIds] 수업 가능한 요가 아이디
@apiBody {int[]} [sigunguIds] 활동가능한 지역 아이디
@apiSuccess (201 or 200) {Number} code 201 or 200
@apiSuccess (201 or 200) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *teacherHandler) Update(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)
	teacherId := ctx.UserValue("teacher_id").(int)

	reqBody := new(request.TeacherUpdateBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.TeacherUpdateParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	// 소유권체크
	if teacherId != reqParam.Id {
		return c.Status(http.StatusForbidden).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	isUpdated, err := h.teacherUsecase.Update(ctx, reqBody, reqParam.Id, userId)
	if err != nil {
		return utils.NewError(c, err)
	}

	httpCode := http.StatusOK
	if !isUpdated {
		httpCode = http.StatusCreated
	}

	return c.Status(httpCode).JSON(ex.Response{
		Code:    httpCode,
		Message: "",
	})
}

// 선생님 부분 수정
/**
@api {patch} /teacher/:id 선생님 부분 수정
@apiName patchTeacher
@apiVersion 1.0.0
@apiGroup teacher
@apiDescription 선생님 부분 수정
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 선생님 아이디
@apiBody {Object} [info]
@apiBody {String} [info.name] 선생님 이름
@apiBody {String} [info.profileImageUrl] 대표 이미지 주소
@apiBody {Number} [info.age] 나이
@apiBody {String} [info.introduce] 간단한 자기소개
@apiBody {Boolean} [info.isProfileOpen] 프로필 공개여부
@apiBody {Object[]} [workExperiences] 근무 경험
@apiBody {Number} [workExperiences.id] 아이디
@apiBody {String} [workExperiences.academyName] 근무지 학원 이름
@apiBody {String} [workExperiences.workStartAt] 근무시작 일시
@apiBody {String} [workExperiences.WorkEndAt] 근무종료 일시
@apiBody {String} [workExperiences.description] 근무에 대한 설명
@apiBody {Object[]} [certifications] 자격증
@apiBody {Number} [certifications.id] 아이디
@apiBody {String} [certifications.agencyName] 취득 기관 이름
@apiBody {String} [certifications.imageUrl] 자격증 이미지 주소
@apiBody {String} [certifications.classStartAt] 수업시작 일시
@apiBody {String} [certifications.classEndAt] 수업종료 일시
@apiBody {String} [certifications.description] 자격증에 대한 설명
@apiBody {Object[]} [yogaRaws] 직접 등록한 요가 (검색 시스템에 나오지 않는 경우 회원이 취급하는 요가 종류 등록)
@apiBody {Number} [yogaRaws.id] 아이디
@apiBody {String} [yogaRaws.name] 요가 이름
@apiBody {int[]} [yogaIds] 수업 가능한 요가 아이디
@apiBody {int[]} [sigunguIds] 활동가능한 지역 아이디
@apiSuccess (201 or 200) {Number} code 201 or 200
@apiSuccess (201 or 200) {String} message ""
@apiError JsonMissing <code>400</code> code: 3000
@apiError ValidationError <code>400</code> code: 2xxx
@apiError YogaDoseNotExist <code>409</code> code: 4007
@apiError SigunguDoseNotExist <code>409</code> code: 4009
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *teacherHandler) Patch(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)
	teacherId := ctx.UserValue("teacher_id").(int)

	reqBody := new(request.TeacherPatchBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}
	if err := h.Validator.ValidateStruct(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewInvalidInputError(err))
	}

	reqParam := new(request.TeacherPatchParam)
	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrJsonMissing, nil))
	}

	// 소유권체크
	if teacherId != reqParam.Id {
		return c.Status(http.StatusForbidden).JSON(ex.NewHttpError(ex.ErrOnlyOwnUser, nil))
	}

	isUpdated, err := h.teacherUsecase.Patch(ctx, reqBody, reqParam.Id, userId)
	if err != nil {
		return utils.NewError(c, err)
	}
	httpCode := http.StatusOK
	if !isUpdated {
		httpCode = http.StatusCreated
	}

	return c.Status(httpCode).JSON(ex.Response{
		Code:    httpCode,
		Message: "",
	})
}

// 선생님 상세 조회
/**
@api {get} /teacher/:id 선생님 상세 조회
@apiName getTeacher
@apiVersion 1.0.0
@apiGroup teacher
@apiDescription 선생님 상세 조회
@apiHeader Authorization accessToken (Bearer)
@apiParam {Number} id 선생님 아이디
@apiSuccess {Number} code 200
@apiSuccess {String} message ""
@apiSuccess {Object} result
@apiSuccess {Number} result.id 아이디
@apiSuccess {String} result.name 섯냉님 이름
@apiSuccess {String} [result.age] 나이
@apiSuccess {String} [result.introduce] 자기소개
@apiSuccess {String} [result.profileImageUrl] 프로필 이미지 주소
@apiSuccess {String} [result.isProfileOpen] 프로필 공개 여부
@apiSuccess {Object[]} [result.yoga] 요가
@apiSuccess {Number} result.yoga.index 요가 순서 번호
@apiSuccess {Number} result.yoga.id 요가 아이디
@apiSuccess {String} result.yoga.nameKor 요가 한글 이름
@apiSuccess {String} result.yoga.isReference 정식적으로 등록됐는지 여부
@apiSuccess {Object[]} [result.possibleWorkSigungu] 근무 가능한 지역
@apiSuccess {Number} result.possibleWorkSigungu.id 근무 가능한 지역 아이디
@apiSuccess {String} result.possibleWorkSigungu.nameKor 시군구 이름
@apiSuccess {Object[]} [result.certifications] 자격증
@apiSuccess {Number} result.certifications.id 자격증 아이디
@apiSuccess {string} result.certifications.agencyName 자격증 기관명
@apiSuccess {string} [result.certifications.imageUrl] 이미지 주소
@apiSuccess {string} result.certifications.classStartAt 수업 시작일
@apiSuccess {string} [result.certifications.classEndAt]  수업 종료일
@apiSuccess {string} [result.certifications.description] 설명
@apiSuccess {string} result.certifications.createdAt 생성일시
@apiSuccess {string} result.certifications.updatedAt 업데이트일시
@apiSuccess {Object[]} [result.workExperiences] 근무경험
@apiSuccess {Number} result.workExperiences.id 근무경험 아이디
@apiSuccess {string} result.certifications.academyName 학원 이름
@apiSuccess {string} result.certifications.workStartAt 수업 시작일
@apiSuccess {string} [result.certifications.workEndAt]  수업 종료일
@apiSuccess {string} [result.certifications.description] 설명
@apiSuccess {string} result.workExperiences.createdAt 생성일시
@apiSuccess {string} result.workExperiences.updatedAt 업데이트일시
@apiSuccess {String} result.createdAt 생성일시
@apiSuccess {String} result.updatedAt 업데이트일시
@apiError ParamsMissing <code>400</code> code: 3002
@apiError TeacherNotFound <code>404</code> code: 5005
@apiError InternalServerError <code>500</code> code: 500
*/
func (h *teacherHandler) Get(c *fiber.Ctx) error {
	ctx := c.Context()
	userId := ctx.UserValue("user_id").(int)

	reqParam := new(request.AcademyDetailParam)

	if err := c.ParamsParser(reqParam); err != nil {
		return c.Status(http.StatusBadRequest).JSON(ex.NewHttpError(ex.ErrParamsMissing, err.Error()))
	}

	teacher, err := h.teacherUsecase.Get(ctx, reqParam.Id, userId)
	if err != nil {
		return utils.NewError(c, err)
	}

	resp := response.NewTeacherResponse(teacher)

	return c.Status(http.StatusOK).JSON(ex.ResponseWithData{
		Code:    http.StatusOK,
		Message: "",
		Result:  resp,
	})
}
