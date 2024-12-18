// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UUIDCreator is an autogenerated mock type for the UUIDCreator type
type UUIDCreator struct {
	mock.Mock
}

// New provides a mock function with given fields:
func (_m *UUIDCreator) New() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for New")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewUUIDCreator creates a new instance of UUIDCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUUIDCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *UUIDCreator {
	mock := &UUIDCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
