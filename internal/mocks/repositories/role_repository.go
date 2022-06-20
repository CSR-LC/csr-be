// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"
)

// RoleRepository is an autogenerated mock type for the RoleRepository type
type RoleRepository struct {
	mock.Mock
}

// GetRoles provides a mock function with given fields: ctx
func (_m *RoleRepository) GetRoles(ctx context.Context) ([]*ent.Role, error) {
	ret := _m.Called(ctx)

	var r0 []*ent.Role
	if rf, ok := ret.Get(0).(func(context.Context) []*ent.Role); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRoleRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoleRepository creates a new instance of RoleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoleRepository(t mockConstructorTestingTNewRoleRepository) *RoleRepository {
	mock := &RoleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
