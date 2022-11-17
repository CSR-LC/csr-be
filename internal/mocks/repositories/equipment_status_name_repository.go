// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"
)

// EquipmentStatusNameRepository is an autogenerated mock type for the EquipmentStatusNameRepository type
type EquipmentStatusNameRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, name
func (_m *EquipmentStatusNameRepository) Create(ctx context.Context, name string) (*ent.EquipmentStatusName, error) {
	ret := _m.Called(ctx, name)

	var r0 *ent.EquipmentStatusName
	if rf, ok := ret.Get(0).(func(context.Context, string) *ent.EquipmentStatusName); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.EquipmentStatusName)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *EquipmentStatusNameRepository) Delete(ctx context.Context, id int) (*ent.EquipmentStatusName, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.EquipmentStatusName
	if rf, ok := ret.Get(0).(func(context.Context, int) *ent.EquipmentStatusName); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.EquipmentStatusName)
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

// Get provides a mock function with given fields: ctx, id
func (_m *EquipmentStatusNameRepository) Get(ctx context.Context, id int) (*ent.EquipmentStatusName, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.EquipmentStatusName
	if rf, ok := ret.Get(0).(func(context.Context, int) *ent.EquipmentStatusName); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.EquipmentStatusName)
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

// GetAll provides a mock function with given fields: ctx
func (_m *EquipmentStatusNameRepository) GetAll(ctx context.Context) ([]*ent.EquipmentStatusName, error) {
	ret := _m.Called(ctx)

	var r0 []*ent.EquipmentStatusName
	if rf, ok := ret.Get(0).(func(context.Context) []*ent.EquipmentStatusName); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.EquipmentStatusName)
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

type mockConstructorTestingTNewEquipmentStatusNameRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewEquipmentStatusNameRepository creates a new instance of EquipmentStatusNameRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEquipmentStatusNameRepository(t mockConstructorTestingTNewEquipmentStatusNameRepository) *EquipmentStatusNameRepository {
	mock := &EquipmentStatusNameRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}