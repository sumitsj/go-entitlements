// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/sumitsj/go-entitlements/gorm/model"

	respository "github.com/sumitsj/go-entitlements/gorm/respository"

	uuid "github.com/google/uuid"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CommitTransaction provides a mock function with given fields:
func (_m *Repository) CommitTransaction() {
	_m.Called()
}

// CreateEntitlementHistory provides a mock function with given fields: ctx, entitlementHistory
func (_m *Repository) CreateEntitlementHistory(ctx context.Context, entitlementHistory *model.UserEntitlementHistory) error {
	ret := _m.Called(ctx, entitlementHistory)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntitlementHistory) error); ok {
		r0 = rf(ctx, entitlementHistory)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetEntitlementsBy provides a mock function with given fields: ctx, userId
func (_m *Repository) GetEntitlementsBy(ctx context.Context, userId uuid.UUID) ([]model.UserEntitlement, error) {
	ret := _m.Called(ctx, userId)

	var r0 []model.UserEntitlement
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) []model.UserEntitlement); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserEntitlement)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RollbackTransaction provides a mock function with given fields:
func (_m *Repository) RollbackTransaction() {
	_m.Called()
}

// UpsertEntitlement provides a mock function with given fields: ctx, entitlement
func (_m *Repository) UpsertEntitlement(ctx context.Context, entitlement *model.UserEntitlement) error {
	ret := _m.Called(ctx, entitlement)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntitlement) error); ok {
		r0 = rf(ctx, entitlement)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithTransaction provides a mock function with given fields:
func (_m *Repository) WithTransaction() respository.Repository {
	ret := _m.Called()

	var r0 respository.Repository
	if rf, ok := ret.Get(0).(func() respository.Repository); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(respository.Repository)
		}
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
