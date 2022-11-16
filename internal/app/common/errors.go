package common

import (
	"fmt"
	"net/http"
)

// ------------------- Error Code And Message -------------------
const (
	ErrInternalError = 500

	// Unprocessable
	ErrJsonMissing        = 1000 // JSON을 입력해주세요
	ErrQueryStringMissing = 1001 // QueryString을 입력해주세요
	ErrParamsMissing      = 1002

	// 벨리데이션 에러 2000 ~
	ErrReqeustsInvalid = 2000
	ErrPasswordInvalid = 2001
	ErrEmailInvalid    = 2002
	ErrUrlInvalid      = 2003
	ErrPhoneNumInvalid = 2004

	// 3000 ~ BadReqeust
	ErrImageExtensionUnavailable            = 3001
	ErrFormDataKeyUnavailable               = 3002
	ErrAuthorizationHeaderFormatUnavailable = 3003
	ErrRandomKeyForEmailVerfiyUnavailable   = 3004
	ErrTokenInvalid                         = 3005
	ErrBusinessCodeInvalid                  = 3006

	// 4000 ~ Conflict
	ErrUserEmailAlreadyExist   = 4001
	ErrUserEmailAlreadyVerfied = 4002
	ErrUserTypeAlreadyRegisted = 4003

	// 5000 ~ NotFound
	ErrUserNotFound      = 5001
	ErrUserEmailNotFound = 5002
	ErrAcademyNotFound   = 5003

	// 6000 ~ 401 Authentication UnAuthorization
	ErrUserEmailUnauthorization = 6001
	ErrTokenExpired             = 6002
	ErrEmailForVerifyExpired    = 6003

	// 7000 ~ Forbidden
)

func ErrorText(code int) string {
	switch code {
	case ErrInternalError:
		return "일시적인 에러가 발생했습니다."
	// 1000 ~
	case ErrJsonMissing:
		return "JSON을 입력해주세요."
	case ErrQueryStringMissing:
		return "쿼리스트링을 입력해주세요."
	case ErrParamsMissing:
		return "파라메터를 입력해주세요."
	// 2000 ~
	case ErrReqeustsInvalid:
		return "유효하지 않은 요청값들이 존재합니다."
	case ErrPasswordInvalid:
		return "유효하지 않은 패스워드입니다."
	case ErrPhoneNumInvalid:
		return "유효하지 않은 휴대폰번호입니다."

	// 3000 ~
	case ErrImageExtensionUnavailable:
		return "이미지 파일이 아닙니다."
	case ErrFormDataKeyUnavailable:
		return "폼 데이터 Key를 확인해주세요."
	case ErrAuthorizationHeaderFormatUnavailable:
		return "Authorization 헤더 형식을 확인해주세요."
	case ErrRandomKeyForEmailVerfiyUnavailable:
		return "사용할 수 없는 인증 키입니다."
	case ErrTokenInvalid:
		return "유효하지 않은 토큰입니다."
	case ErrBusinessCodeInvalid:
		return "유효하지 않은 사업자 번호입니다."

	// 4000 ~ Conflict
	case ErrUserEmailAlreadyExist:
		return "이미 존재하는 이메일입니다."
	case ErrUserEmailAlreadyVerfied:
		return "이미 이메일 인증이 완료됐습니다."
	case ErrUserTypeAlreadyRegisted:
		return "이미 회원 유형이 등록됐습니다."

	// 5000 ~
	case ErrUserNotFound:
		return "존재하지 않는 유저입니다."
	case ErrAcademyNotFound:
		return "존재하지 않는 학원입니다."

	// 6000 ~
	case ErrUserEmailUnauthorization:
		return "유저의 이메일 인증이 필요합니다."

	case ErrTokenExpired:
		return "토큰이 만료되었습니다."

	case ErrEmailForVerifyExpired:
		return "이메일 인증 시간이 만료되었습니다."

	default:
		return "일시적인 에러가 발생했습니다."
	}
}

// ------------------- Make Error Method -------------------
type HttpErr interface {
	HttpStatusCode() int
	ErrorCode() int
	Error() string
	Details() interface{}
}

type HttpError struct {
	ErrHttpCode int         `json:"httpCode,omitempty"`
	ErrCode     int         `json:"code"`
	ErrMessage  string      `json:"message"`
	ErrDetails  interface{} `json:"details"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - details: %v", e.ErrCode, e.ErrMessage, e.ErrDetails)
}

// Error status
func (e HttpError) HttpStatusCode() int {
	return e.ErrHttpCode
}

func (e HttpError) ErrorCode() int {
	return e.ErrCode
}

// HttpError Details
func (e HttpError) Details() interface{} {
	return e.ErrDetails
}

func NewHttpError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrCode:    ErrorCode,
		ErrMessage: ErrorText(ErrorCode),
		ErrDetails: details,
	}
}

// 400
func NewBadRequestError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusBadRequest,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

// 401
func NewAuthenticationFailedError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: 401,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

func NewUnauthorizedError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusUnauthorized,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

func NewForbiddenError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusForbidden,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

func NewConflictError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusConflict,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

// New Not Found Error
func NewNotFoundError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusNotFound,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

// New Unprocessable Entity Error
func NewUnprocessableEntityError(ErrorCode int, details interface{}) HttpErr {
	return HttpError{
		ErrHttpCode: http.StatusUnprocessableEntity,
		ErrCode:     ErrorCode,
		ErrMessage:  ErrorText(ErrorCode),
		ErrDetails:  details,
	}
}

// New Internal Server Error
func NewInternalServerError() HttpErr {
	// TODO : loging
	return HttpError{
		ErrCode:    ErrInternalError,
		ErrMessage: ErrorText(ErrInternalError),
		ErrDetails: nil,
	}
}

func ParseHttpError(err error) (int, interface{}) {
	if httpErr, ok := err.(HttpErr); ok {
		return httpErr.HttpStatusCode(), NewHttpError(httpErr.ErrorCode(), httpErr.Details())
	}
	return http.StatusInternalServerError, NewInternalServerError()
}
