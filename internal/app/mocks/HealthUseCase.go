// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// HealthUseCase is an autogenerated mock type for the HealthUseCase type
type HealthUseCase struct {
	mock.Mock
}

type mockConstructorTestingTNewHealthUseCase interface {
	mock.TestingT
	Cleanup(func())
}

// NewHealthUseCase creates a new instance of HealthUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewHealthUseCase(t mockConstructorTestingTNewHealthUseCase) *HealthUseCase {
	mock := &HealthUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
