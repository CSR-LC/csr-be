// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	ent "git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	mock "github.com/stretchr/testify/mock"

	models "git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// PhotoRepository is an autogenerated mock type for the PhotoRepository type
type PhotoRepository struct {
	mock.Mock
}

// CreatePhoto provides a mock function with given fields: ctx, p
func (_m *PhotoRepository) CreatePhoto(ctx context.Context, p models.Photo) (*ent.Photo, error) {
	ret := _m.Called(ctx, p)

	var r0 *ent.Photo
	if rf, ok := ret.Get(0).(func(context.Context, models.Photo) *ent.Photo); ok {
		r0 = rf(ctx, p)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Photo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Photo) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePhotoByID provides a mock function with given fields: ctx, id
func (_m *PhotoRepository) DeletePhotoByID(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PhotoByID provides a mock function with given fields: ctx, id
func (_m *PhotoRepository) PhotoByID(ctx context.Context, id string) (*ent.Photo, error) {
	ret := _m.Called(ctx, id)

	var r0 *ent.Photo
	if rf, ok := ret.Get(0).(func(context.Context, string) *ent.Photo); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ent.Photo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPhotoRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPhotoRepository creates a new instance of PhotoRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPhotoRepository(t mockConstructorTestingTNewPhotoRepository) *PhotoRepository {
	mock := &PhotoRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}