// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	validatorx "onthemat/pkg/validatorx"

	mock "github.com/stretchr/testify/mock"
)

// Validator is an autogenerated mock type for the Validator type
type Validator struct {
	mock.Mock
}

// ValidateStruct provides a mock function with given fields: request
func (_m *Validator) ValidateStruct(request interface{}) []*validatorx.ErrorResponse {
	ret := _m.Called(request)

	var r0 []*validatorx.ErrorResponse
	if rf, ok := ret.Get(0).(func(interface{}) []*validatorx.ErrorResponse); ok {
		r0 = rf(request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*validatorx.ErrorResponse)
		}
	}

	return r0
}

type mockConstructorTestingTNewValidator interface {
	mock.TestingT
	Cleanup(func())
}

// NewValidator creates a new instance of Validator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewValidator(t mockConstructorTestingTNewValidator) *Validator {
	mock := &Validator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}