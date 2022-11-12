// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	ent "onthemat/pkg/ent"

	mock "github.com/stretchr/testify/mock"
)

// TeacherWorkExperience is an autogenerated mock type for the TeacherWorkExperience type
type TeacherWorkExperience struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, value, teacherId
func (_m *TeacherWorkExperience) Create(ctx context.Context, value *ent.TeacherWorkExperience, teacherId int) error {
	ret := _m.Called(ctx, value, teacherId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ent.TeacherWorkExperience, int) error); ok {
		r0 = rf(ctx, value, teacherId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateMany provides a mock function with given fields: ctx, value, teacherId
func (_m *TeacherWorkExperience) CreateMany(ctx context.Context, value []*ent.TeacherWorkExperience, teacherId int) error {
	ret := _m.Called(ctx, value, teacherId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ent.TeacherWorkExperience, int) error); ok {
		r0 = rf(ctx, value, teacherId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, workExperienceId
func (_m *TeacherWorkExperience) Delete(ctx context.Context, workExperienceId int) error {
	ret := _m.Called(ctx, workExperienceId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, workExperienceId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, workExperienceId
func (_m *TeacherWorkExperience) Get(ctx context.Context, workExperienceId int) (*ent.TeacherWorkExperience, error) {
	ret := _m.Called(ctx, workExperienceId)

	var r0 *ent.TeacherWorkExperience
	if rf, ok := ret.Get(0).(func(context.Context, int) *ent.TeacherWorkExperience); ok {
		r0 = rf(ctx, workExperienceId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.TeacherWorkExperience)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, workExperienceId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByTeacherID provides a mock function with given fields: ctx, teacherId
func (_m *TeacherWorkExperience) ListByTeacherID(ctx context.Context, teacherId int) ([]*ent.TeacherWorkExperience, error) {
	ret := _m.Called(ctx, teacherId)

	var r0 []*ent.TeacherWorkExperience
	if rf, ok := ret.Get(0).(func(context.Context, int) []*ent.TeacherWorkExperience); ok {
		r0 = rf(ctx, teacherId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.TeacherWorkExperience)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, teacherId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, value
func (_m *TeacherWorkExperience) Update(ctx context.Context, value *ent.TeacherWorkExperience) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ent.TeacherWorkExperience) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTeacherWorkExperience interface {
	mock.TestingT
	Cleanup(func())
}

// NewTeacherWorkExperience creates a new instance of TeacherWorkExperience. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTeacherWorkExperience(t mockConstructorTestingTNewTeacherWorkExperience) *TeacherWorkExperience {
	mock := &TeacherWorkExperience{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
