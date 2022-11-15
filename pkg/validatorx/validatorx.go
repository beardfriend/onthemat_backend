package validatorx

import (
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
}

func (v *validator) AddPasswordAtLeastOneCharNumValidation(tag string) ValidatorSetter {
	if err := v.validator.RegisterValidation(tag, validatePassword); err != nil {
		panic(err)
	}
	return v
}

func validatePassword(fl goValidator.FieldLevel) bool {
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
	FailedField string
	Tag         string
	Value       string
}

func (v *validator) ValidateStruct(request interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := v.validator.Struct(request)
	if err != nil {
		for _, err := range err.(goValidator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
