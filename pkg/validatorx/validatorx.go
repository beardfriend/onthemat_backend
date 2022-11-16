package validatorx

import (
	"reflect"
	"strings"

	goValidator "github.com/go-playground/validator/v10"
)

// ------------------- Validator -------------------
type Validator interface {
	ValidateStruct(request interface{}) []*ErrorResponse
}

// ------------------- New -------------------
type validator struct {
	validator *goValidator.Validate
}

func NewValidatorx() ValidatorSetter {
	validate := goValidator.New()

	return &validator{
		validator: validate,
	}
}

// ------------------- Setter -------------------

type ValidatorSetter interface {
	Init() Validator
	AddPasswordAtLeastOneCharNumValidation(tag string) ValidatorSetter
	AddPhoneNumValidation(tagName string) ValidatorSetter
	AddUrlValidation(tagName string) ValidatorSetter
	SetExtractTagName() ValidatorSetter
}

func (v *validator) AddPasswordAtLeastOneCharNumValidation(tagName string) ValidatorSetter {
	if err := v.validator.RegisterValidation(tagName, validatePasswordAtLeastOneCharNum); err != nil {
		panic(err)
	}
	return v
}

func (v *validator) AddPhoneNumValidation(tagName string) ValidatorSetter {
	if err := v.validator.RegisterValidation(tagName, validatePhoneNum); err != nil {
		panic(err)
	}
	return v
}

func (v *validator) AddUrlValidation(tagName string) ValidatorSetter {
	if err := v.validator.RegisterValidation(tagName, validateUrl); err != nil {
		panic(err)
	}
	return v
}

func (v *validator) SetExtractTagName() ValidatorSetter {
	v.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		queryName := strings.SplitN(field.Tag.Get("query"), ",", 2)[0]
		paramName := strings.SplitN(field.Tag.Get("param"), ",", 2)[0]

		if jsonName != "" && jsonName != "-" {
			return jsonName
		}

		if queryName != "" && queryName != "-" {
			return queryName
		}

		if paramName != "" && paramName != "-" {
			return paramName
		}

		return ""
	})
	return v
}

func validatePhoneNum(fl goValidator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return ValidateRegex(PhoneRegex, value)
	}
	return true
}

func validateUrl(fl goValidator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return ValidateRegex(UrlRegex, value)
	}
	return true
}

func validatePasswordAtLeastOneCharNum(fl goValidator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		return ValidateRegex(ForbiddenSpecialCharRegex, value) && ValidateRegex(AtLeastOneCharOneNumRegex, value)
	}
	return true
}

func (v *validator) Init() Validator {
	return v
}

// ------------------- Validator Method -------------------
type ErrorResponse struct {
	FailedField        string
	FailedFieldTagName string
	Tag                string
	Value              string
}

func (v *validator) ValidateStruct(request interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := v.validator.Struct(request)
	if err != nil {
		for _, err := range err.(goValidator.ValidationErrors) {
			var element ErrorResponse

			element.FailedFieldTagName = err.Field()
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
