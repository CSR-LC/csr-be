// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"
)

// OrderStatusNameRepository is an autogenerated mock type for the OrderStatusNameRepository type
type OrderStatusNameRepository struct {
	mock.Mock
}

// ListOfOrderStatusNames provides a mock function with given fields: ctx
func (_m *OrderStatusNameRepository) ListOfOrderStatusNames(ctx context.Context) ([]*ent.OrderStatusName, error) {
	ret := _m.Called(ctx)

	var r0 []*ent.OrderStatusName
	if rf, ok := ret.Get(0).(func(context.Context) []*ent.OrderStatusName); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.OrderStatusName)
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

type mockConstructorTestingTNewOrderStatusNameRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewOrderStatusNameRepository creates a new instance of OrderStatusNameRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOrderStatusNameRepository(t mockConstructorTestingTNewOrderStatusNameRepository) *OrderStatusNameRepository {
	mock := &OrderStatusNameRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}