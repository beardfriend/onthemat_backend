package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"onthemat/pkg/validatorx"
)

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrAuthenticationFailed = errors.New("authentication vailed")
	ErrBadRequest           = errors.New("bad request")
	ErrNotFound             = errors.New("not found")
	ErrUnprocessableEntity  = errors.New("unprocessable entity")
	ErrInternalServerError  = errors.New("internal server error")
)

type HttpErr interface {
	Status() int
	Error() string
	Details() interface{}
}

type HttpError struct {
	ErrCode    int         `json:"code"`
	ErrMessage string      `json:"message"`
	ErrDetails interface{} `json:"details"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - details: %v", e.ErrCode, e.ErrMessage, e.ErrDetails)
}

// Error status
func (e HttpError) Status() int {
	return e.ErrCode
}

// HttpError Details
func (e HttpError) Details() interface{} {
	return e.ErrDetails
}

func NewHttpError(status int, err string, details interface{}) HttpErr {
	return HttpError{
		ErrCode:    status,
		ErrMessage: err,
		ErrDetails: details,
	}
}

func NewAuthenticationFailedError(details interface{}) HttpErr {
	return HttpError{
		ErrCode:    401,
		ErrMessage: ErrAuthenticationFailed.Error(),
		ErrDetails: details,
	}
}

func NewUnauthorizedError(details interface{}) HttpErr {
	return HttpError{
		ErrCode:    http.StatusUnauthorized,
		ErrMessage: ErrUnauthorized.Error(),
		ErrDetails: details,
	}
}

// New Unprocessable Entity Error
func NewUnprocessableEntityError(details interface{}) HttpErr {
	return HttpError{
		ErrCode:    http.StatusUnprocessableEntity,
		ErrMessage: ErrUnprocessableEntity.Error(),
		ErrDetails: details,
	}
}

// New Internal Server Error
func NewInternalServerError(details interface{}) HttpErr {
	// TODO : loging
	return HttpError{
		ErrCode:    http.StatusInternalServerError,
		ErrMessage: ErrInternalServerError.Error(),
		ErrDetails: nil,
	}
}

// New Invalid Input Error - Validation
func NewInvalidInputError(errs []*validatorx.ErrorResponse) HttpErr {
	var errors []interface{}
	for _, field := range errs {
		splited := strings.Split(field.FailedField, ".")
		fieldName := splited[1]
		errors = append(errors, map[string]interface{}{
			fieldName: field.Tag + field.Value,
		})

	}

	return HttpError{
		ErrCode:    http.StatusBadRequest,
		ErrMessage: ErrBadRequest.Error(),
		ErrDetails: errors,
	}
}
