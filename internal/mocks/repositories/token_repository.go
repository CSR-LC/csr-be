// Code generated by mockery v2.13.1. DO NOT EDIT.

package repositories

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TokenRepository is an autogenerated mock type for the TokenRepository type
type TokenRepository struct {
	mock.Mock
}

// CreateTokens provides a mock function with given fields: ctx, ownerID, accessToken, refreshToken
func (_m *TokenRepository) CreateTokens(ctx context.Context, ownerID int, accessToken string, refreshToken string) error {
	ret := _m.Called(ctx, ownerID, accessToken, refreshToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string, string) error); ok {
		r0 = rf(ctx, ownerID, accessToken, refreshToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTokenRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTokenRepository creates a new instance of TokenRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTokenRepository(t mockConstructorTestingTNewTokenRepository) *TokenRepository {
	mock := &TokenRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
