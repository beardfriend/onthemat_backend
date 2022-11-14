// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	ent "onthemat/pkg/ent"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, user
func (_m *UserRepository) Create(ctx context.Context, user *ent.User) (*ent.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, *ent.User) *ent.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ent.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepository) FindByEmail(ctx context.Context, email string) (bool, error) {
	ret := _m.Called(ctx, email)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *UserRepository) Get(ctx context.Context, id int) (*ent.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, int) *ent.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	ret := _m.Called(ctx, email)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *ent.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmailPassword provides a mock function with given fields: ctx, u
func (_m *UserRepository) GetByEmailPassword(ctx context.Context, u *ent.User) (*ent.User, error) {
	ret := _m.Called(ctx, u)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, *ent.User) *ent.User); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ent.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBySocialKey provides a mock function with given fields: ctx, u
func (_m *UserRepository) GetBySocialKey(ctx context.Context, u *ent.User) (*ent.User, error) {
	ret := _m.Called(ctx, u)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, *ent.User) *ent.User); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ent.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, user
func (_m *UserRepository) Update(ctx context.Context, user *ent.User) (*ent.User, error) {
	ret := _m.Called(ctx, user)

	var r0 *ent.User
	if rf, ok := ret.Get(0).(func(context.Context, *ent.User) *ent.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *ent.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEmailVerifeid provides a mock function with given fields: ctx, userId
func (_m *UserRepository) UpdateEmailVerifeid(ctx context.Context, userId int) error {
	ret := _m.Called(ctx, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTempPassword provides a mock function with given fields: ctx, u
func (_m *UserRepository) UpdateTempPassword(ctx context.Context, u *ent.User) error {
	ret := _m.Called(ctx, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ent.User) error); ok {
		r0 = rf(ctx, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserRepository(t mockConstructorTestingTNewUserRepository) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
