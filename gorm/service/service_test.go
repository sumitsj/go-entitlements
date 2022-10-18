package service_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sumitsj/go-entitlements/gorm/model"
	"github.com/sumitsj/go-entitlements/gorm/respository/mocks"
	"github.com/sumitsj/go-entitlements/gorm/service"
	"gorm.io/gorm"
	"testing"
)

func Test_NewService(t *testing.T) {
	repository := mocks.NewRepository(t)
	service := service.NewService(repository)
	assert.NotNil(t, service)
}

func Test_GetEntitlementsBy(t *testing.T) {
	ctx := context.TODO()
	userId := uuid.New()

	t.Run("should process successfully", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repository.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{{}}, nil)

		service := service.NewService(repository)

		entitlements, err := service.GetEntitlementsBy(ctx, userId)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(entitlements))
	})

	t.Run("should return error on repository error", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repository.On("GetEntitlementsBy", ctx, userId).
			Return([]model.UserEntitlement{}, errors.New("some error"))

		service := service.NewService(repository)

		_, err := service.GetEntitlementsBy(ctx, userId)

		assert.Error(t, err)
	})
}

func Test_UpsertEntitlement(t *testing.T) {
	ctx := context.TODO()
	userId := uuid.New()
	entitlementName := "test_entitlement"
	reason := "disable"

	t.Run("should process successfully when entitlement already exists", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repositoryWithTransaction := mocks.NewRepository(t)
		repository.On("WithTransaction").Return(repositoryWithTransaction)
		repositoryWithTransaction.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{{
			Model:  gorm.Model{ID: 1},
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}}, nil)
		repositoryWithTransaction.On("UpsertEntitlement", ctx, mock.MatchedBy(entitlement(userId, entitlementName, false))).
			Return(nil)
		repositoryWithTransaction.On("CreateEntitlementHistory", ctx, mock.MatchedBy(entitlementHistory(1, reason, true, false))).
			Return(nil)
		repositoryWithTransaction.On("CommitTransaction")

		service := service.NewService(repository)

		err := service.UpsertEntitlement(ctx, &model.UserEntitlement{
			Model:  gorm.Model{ID: 1},
			UserId: userId,
			Name:   entitlementName,
			Value:  false,
		}, reason)

		assert.NoError(t, err)
	})

	t.Run("should process successfully when entitlement does not exist", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repositoryWithTransaction := mocks.NewRepository(t)
		repository.On("WithTransaction").Return(repositoryWithTransaction)
		repositoryWithTransaction.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, nil)
		repositoryWithTransaction.On("UpsertEntitlement", ctx, mock.MatchedBy(entitlement(userId, entitlementName, true))).
			Return(nil)
		repositoryWithTransaction.On("CreateEntitlementHistory", ctx, mock.MatchedBy(entitlementHistory(0, reason, false, true))).
			Return(nil)
		repositoryWithTransaction.On("CommitTransaction")

		service := service.NewService(repository)

		err := service.UpsertEntitlement(ctx, &model.UserEntitlement{
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}, reason)

		assert.NoError(t, err)
	})

	t.Run("should return error on error while getting entitlement", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repositoryWithTransaction := mocks.NewRepository(t)
		repository.On("WithTransaction").Return(repositoryWithTransaction)
		repositoryWithTransaction.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, errors.New("some error"))
		repositoryWithTransaction.On("RollbackTransaction")

		service := service.NewService(repository)

		err := service.UpsertEntitlement(ctx, &model.UserEntitlement{
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}, reason)

		assert.Error(t, err)
	})

	t.Run("should return error on error while updating entitlement", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repositoryWithTransaction := mocks.NewRepository(t)
		repository.On("WithTransaction").Return(repositoryWithTransaction)
		repositoryWithTransaction.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, nil)
		repositoryWithTransaction.On("UpsertEntitlement", ctx, mock.MatchedBy(entitlement(userId, entitlementName, true))).
			Return(errors.New("some error"))
		repositoryWithTransaction.On("RollbackTransaction")

		service := service.NewService(repository)

		err := service.UpsertEntitlement(ctx, &model.UserEntitlement{
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}, reason)

		assert.Error(t, err)
	})

	t.Run("should return error on error while saving entitlement history", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repositoryWithTransaction := mocks.NewRepository(t)
		repository.On("WithTransaction").Return(repositoryWithTransaction)
		repositoryWithTransaction.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, nil)
		repositoryWithTransaction.On("UpsertEntitlement", ctx, mock.MatchedBy(entitlement(userId, entitlementName, true))).
			Return(nil)
		repositoryWithTransaction.On("CreateEntitlementHistory", ctx, mock.MatchedBy(entitlementHistory(0, reason, false, true))).
			Return(errors.New("some error"))
		repositoryWithTransaction.On("RollbackTransaction")

		service := service.NewService(repository)

		err := service.UpsertEntitlement(ctx, &model.UserEntitlement{
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}, reason)

		assert.Error(t, err)
	})
}

func Test_service_IsEntitlementEnabled(t *testing.T) {
	ctx := context.TODO()
	userId := uuid.New()
	entitlementName := "test_entitlement"

	t.Run("should return true as entitlement is enabled", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repository.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{{
			UserId: userId,
			Name:   entitlementName,
			Value:  true,
		}}, nil)

		service := service.NewService(repository)

		isEnabled, err := service.IsEntitlementEnabled(ctx, userId, entitlementName)

		assert.NoError(t, err)
		assert.True(t, isEnabled)
	})

	t.Run("should return false if entitlement does not exist", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repository.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, nil)

		service := service.NewService(repository)

		isEnabled, err := service.IsEntitlementEnabled(ctx, userId, entitlementName)

		assert.NoError(t, err)
		assert.False(t, isEnabled)
	})

	t.Run("should return error on error while getting entitlement", func(t *testing.T) {
		repository := mocks.NewRepository(t)
		repository.On("GetEntitlementsBy", ctx, userId).Return([]model.UserEntitlement{}, errors.New("some error"))

		service := service.NewService(repository)

		_, err := service.IsEntitlementEnabled(ctx, userId, entitlementName)

		assert.Error(t, err)
	})
}

func entitlement(userId uuid.UUID, name string, value bool) func(entitlement *model.UserEntitlement) bool {
	return func(entitlement *model.UserEntitlement) bool {
		if entitlement.UserId != userId {
			fmt.Printf("\nEntitlement.UserId mis-match. Expected: %s, Actual: %s", userId, entitlement.UserId)
			return false
		}

		if entitlement.Name != name {
			fmt.Printf("\nEntitlement.Name mis-match. Expected: %s, Actual: %s", name, entitlement.Name)
			return false
		}

		if entitlement.Value != value {
			fmt.Printf("\nEntitlement.Value mis-match. Expected: %v, Actual: %v", value, entitlement.Value)
			return false
		}

		return true
	}
}

func entitlementHistory(userEntitlementId uint, reason string, oldValue bool, newValue bool) func(entitlement *model.UserEntitlementHistory) bool {
	return func(entitlementHistory *model.UserEntitlementHistory) bool {
		if entitlementHistory.UserEntitlementId != userEntitlementId {
			fmt.Printf("\nEntitlementHistory.UserEntitlementId mis-match. Expected: %v, Actual: %v", userEntitlementId, entitlementHistory.UserEntitlementId)
			return false
		}

		if entitlementHistory.OldValue != oldValue {
			fmt.Printf("\nEntitlementHistory.OldValue mis-match. Expected: %v, Actual: %v", oldValue, entitlementHistory.OldValue)
			return false
		}

		if entitlementHistory.NewValue != newValue {
			fmt.Printf("\nEntitlementHistory.NewValue mis-match. Expected: %v, Actual: %v", newValue, entitlementHistory.NewValue)
			return false
		}

		if entitlementHistory.Reason != reason {
			fmt.Printf("\nEntitlementHistory.Reason mis-match. Expected: %v, Actual: %v", reason, entitlementHistory.Reason)
			return false
		}

		return true
	}
}
