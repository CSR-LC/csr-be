// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"

	models "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// KindRepository is an autogenerated mock type for the KindRepository type
type KindRepository struct {
	mock.Mock
}

// AllKinds provides a mock function with given fields: ctx, limit, offset, orderBy, orderColumn
func (_m *KindRepository) AllKinds(ctx context.Context, limit int, offset int, orderBy string, orderColumn string) ([]*ent.Kind, error) {
	ret := _m.Called(ctx, limit, offset, orderBy, orderColumn)

	var r0 []*ent.Kind
	if rf, ok := ret.Get(0).(func(context.Context, int, int, string, string) []*ent.Kind); ok {
		r0 = rf(ctx, limit, offset, orderBy, orderColumn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.Kind)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, string, string) error); ok {
		r1 = rf(ctx, limit, offset, orderBy, orderColumn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AllKindsTotal provides a mock function with given fields: ctx
func (_m *KindRepository) AllKindsTotal(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateKind provides a mock function with given fields: ctx, newKind
func (_m *KindRepository) CreateKind(ctx context.Context, newKind models.CreateNewKind) (*ent.Kind, error) {
	ret := _m.Called(ctx, newKind)

	var r0 *ent.Kind
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateNewKind) *ent.Kind); ok {
		r0 = rf(ctx, newKind)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Kind)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.CreateNewKind) error); ok {
		r1 = rf(ctx, newKind)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteKindByID provides a mock function with given fields: ctx, id
func (_m *KindRepository) DeleteKindByID(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// KindByID provides a mock function with given fields: ctx, id
func (_m *KindRepository) KindByID(ctx context.Context, id int) (*ent.Kind, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.Kind
	if rf, ok := ret.Get(0).(func(context.Context, int) *ent.Kind); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Kind)
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

// UpdateKind provides a mock function with given fields: ctx, id, update
func (_m *KindRepository) UpdateKind(ctx context.Context, id int, update models.PatchKind) (*ent.Kind, error) {
	ret := _m.Called(ctx, id, update)

	var r0 *ent.Kind
	if rf, ok := ret.Get(0).(func(context.Context, int, models.PatchKind) *ent.Kind); ok {
		r0 = rf(ctx, id, update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Kind)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, models.PatchKind) error); ok {
		r1 = rf(ctx, id, update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewKindRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewKindRepository creates a new instance of KindRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKindRepository(t mockConstructorTestingTNewKindRepository) *KindRepository {
	mock := &KindRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
