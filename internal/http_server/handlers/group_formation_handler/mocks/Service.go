// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "matchMaker/internal/storage/postgres/repository/models"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetAndRemoveRemainingUsers provides a mock function with given fields: ctx
func (_m *Service) GetAndRemoveRemainingUsers(ctx context.Context) ([]models.User, bool, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAndRemoveRemainingUsers")
	}

	var r0 []models.User
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.User, bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) bool); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetUsersInSearch provides a mock function with given fields: ctx, groupSize
func (_m *Service) GetUsersInSearch(ctx context.Context, groupSize int) ([]models.User, error) {
	ret := _m.Called(ctx, groupSize)

	if len(ret) == 0 {
		panic("no return value specified for GetUsersInSearch")
	}

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]models.User, error)); ok {
		return rf(ctx, groupSize)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.User); ok {
		r0 = rf(ctx, groupSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, groupSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveRemainingUsers provides a mock function with given fields: ctx, users
func (_m *Service) SaveRemainingUsers(ctx context.Context, users []models.User) error {
	ret := _m.Called(ctx, users)

	if len(ret) == 0 {
		panic("no return value specified for SaveRemainingUsers")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []models.User) error); ok {
		r0 = rf(ctx, users)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
