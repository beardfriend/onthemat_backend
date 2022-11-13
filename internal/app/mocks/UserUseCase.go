// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	ent "onthemat/pkg/ent"

	fasthttp "github.com/valyala/fasthttp"

	mock "github.com/stretchr/testify/mock"
)

// UserUseCase is an autogenerated mock type for the UserUseCase type
type UserUseCase struct {
	mock.Mock
}

// GetMe provides a mock function with given fields: ctx, id
func (_m *UserUseCase) GetMe(ctx *fasthttp.RequestCtx, id int) (*ent.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(*fasthttp.RequestCtx, int) *ent.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*fasthttp.RequestCtx, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUseCase creates a new instance of UserUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUseCase(t mockConstructorTestingTNewUserUseCase) *UserUseCase {
	mock := &UserUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}