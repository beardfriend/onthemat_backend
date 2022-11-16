package common

import (
	"fmt"
	"testing"

	"onthemat/pkg/validatorx"
)

func TestNewHttpError(t *testing.T) {
	r := NewHttpError(ErrJsonMissing, nil)
	fmt.Println(r)
}

type TestStruct struct {
	Password string `json:"password" validate:"PassWordAtLeastOneCharOneNum,min=5,max=20"`
	Phone    string `json:"phoneNum" validate:"phoneNumNoDash"`
}

func TestValidationHttpError(t *testing.T) {
	m := validatorx.NewValidatorx().
		AddPasswordAtLeastOneCharNumValidation("PassWordAtLeastOneCharOneNum").
		AddPhoneNumValidation("phoneNumNoDash").SetExtractTagName().Init()

	request := TestStruct{Password: "asdasdasd1", Phone: "010226633"}
	err := m.ValidateStruct(request)
	e := NewInvalidInputError(err)
	fmt.Println(e)
}
