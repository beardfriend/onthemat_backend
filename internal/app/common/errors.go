package common

import (
	"onthemat/pkg/validatorx"
)

const (
	ErrInternalError = 500

	// Unprocessable
	ErrJsonMissing        = 1000 // JSON을 입력해주세요
	ErrQueryStringMissing = 1001 // QueryString을 입력해주세요

	// 벨리데이션 에러 2000 ~
	ErrReqeustsInvalid = 2000
	ErrPasswordInvalid = 2001
	ErrEmailInvalid    = 2002
	ErrUrlInvalid      = 2003
	ErrPhoneNumInvalid = 2004

	// 3000 ~ BadReqeust

	// 4000 ~ Conflict

	// 5000 ~ NotFound

	// 6000 ~ 인증

)

func ErrorText(code int) string {
	switch code {
	// 1000 ~
	case ErrJsonMissing:
		return "JSON을 입력해주세요"
	case ErrQueryStringMissing:
		return "QueryString을 입력해주세요"

	case ErrReqeustsInvalid:
		return "유효하지 않은 요청값들이 존재합니다."
	case ErrPasswordInvalid:
		return "유효하지 않은 패스워드입니다."
	case ErrPhoneNumInvalid:
		return "유효하지 않은 휴대폰번호입니다."
	default:
		return "일시적인 에러가 발생했습니다."
	}
}

type HttpError struct {
	ErrCode    int         `json:"code"`
	ErrMessage string      `json:"message"`
	ErrDetails interface{} `json:"details"`
}

func NewHttpError(ErrorCode int, details interface{}) HttpError {
	return HttpError{
		ErrCode:    ErrorCode,
		ErrMessage: ErrorText(ErrorCode),
		ErrDetails: details,
	}
}

const (
	TagPassword  = "password"
	TagEmail     = "email"
	TagNickName  = "nickname"
	TagPhoneNum  = "phoneNum"
	TagTermAgree = "termAgree"
)

func ErrorValidationCode(name string) int {
	switch name {
	case TagPassword:
		return ErrPasswordInvalid
	case TagEmail:
		return ErrEmailInvalid
	case TagPhoneNum:
		return ErrPhoneNumInvalid

	default:
		return ErrReqeustsInvalid
	}
}

func NewInvalidInputError(errs []*validatorx.ErrorResponse) HttpError {
	errors := make([]interface{}, len(errs))

	for _, field := range errs {
		errors = append(errors, map[string]interface{}{
			field.FailedFieldTagName: field.Tag,
		})
	}

	if len(errs) == 1 {
		return HttpError{
			ErrCode:    ErrorValidationCode(errs[0].FailedFieldTagName),
			ErrMessage: ErrorText(ErrorValidationCode(errs[0].FailedFieldTagName)),
			ErrDetails: errors,
		}
	}

	return HttpError{
		ErrCode:    ErrReqeustsInvalid,
		ErrMessage: ErrorText(ErrReqeustsInvalid),
		ErrDetails: errors,
	}
}
