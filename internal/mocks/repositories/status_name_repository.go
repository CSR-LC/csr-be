// Code generated by mockery v2.9.4. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"
)

// StatusNameRepository is an autogenerated mock type for the StatusNameRepository type
type StatusNameRepository struct {
	mock.Mock
}

// ListOfStatuses provides a mock function with given fields: ctx
func (_m *StatusNameRepository) ListOfStatuses(ctx context.Context) ([]*ent.StatusName, error) {
	ret := _m.Called(ctx)

	var r0 []*ent.StatusName
	if rf, ok := ret.Get(0).(func(context.Context) []*ent.StatusName); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ent.StatusName)
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
