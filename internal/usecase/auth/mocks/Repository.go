// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	dto "auth_service/internal/dto"
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateUserSession provides a mock function with given fields: ctx, req
func (_m *Repository) CreateUserSession(ctx context.Context, req dto.UserSessionRepoCreate) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateUserSession")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.UserSessionRepoCreate) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserSession provides a mock function with given fields: ctx, guid, createdAt
func (_m *Repository) GetUserSession(ctx context.Context, guid string, createdAt time.Time) (dto.UserSessionRepoGet, error) {
	ret := _m.Called(ctx, guid, createdAt)

	if len(ret) == 0 {
		panic("no return value specified for GetUserSession")
	}

	var r0 dto.UserSessionRepoGet
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) (dto.UserSessionRepoGet, error)); ok {
		return rf(ctx, guid, createdAt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) dto.UserSessionRepoGet); ok {
		r0 = rf(ctx, guid, createdAt)
	} else {
		r0 = ret.Get(0).(dto.UserSessionRepoGet)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, time.Time) error); ok {
		r1 = rf(ctx, guid, createdAt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetSessionUsed provides a mock function with given fields: ctx, guid, createdAt
func (_m *Repository) SetSessionUsed(ctx context.Context, guid string, createdAt time.Time) error {
	ret := _m.Called(ctx, guid, createdAt)

	if len(ret) == 0 {
		panic("no return value specified for SetSessionUsed")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Time) error); ok {
		r0 = rf(ctx, guid, createdAt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
