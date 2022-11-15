package validatorx

import (
	"fmt"
	"testing"
)

func TestCustomPassword(t *testing.T) {
	vmodule := NewValidatorx().
		AddPasswordAtLeastOneCharNumValidation("customPass").Init()

	type Request struct {
		Password string `validate:"customPass"`
		Name     string `validate:"required,min=5"`
	}

	req := &Request{
		Password: "As/d0897123!",
		Name:     "asd",
	}
	res := vmodule.ValidateStruct(req)
	fmt.Println(res)
}
