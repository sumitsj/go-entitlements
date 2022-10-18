package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/sumitsj/go-entitlements/gorm/model"
	"github.com/sumitsj/go-entitlements/gorm/respository"
)

type service struct {
	repository respository.Repository
}

type Service interface {
	GetEntitlementsBy(ctx context.Context, portfolioId uuid.UUID) ([]model.UserEntitlement, error)
	UpsertEntitlement(ctx context.Context, entitlement *model.UserEntitlement, reason string) error
	IsEntitlementEnabled(ctx context.Context, portfolioId uuid.UUID, name string) (bool, error)
}

func (s *service) GetEntitlementsBy(ctx context.Context, portfolioId uuid.UUID) ([]model.UserEntitlement, error) {
	return s.repository.GetEntitlementsBy(ctx, portfolioId)
}

func (s *service) UpsertEntitlement(ctx context.Context, entitlement *model.UserEntitlement, reason string) error {
	tx := s.repository.WithTransaction()

	entitlements, err := tx.GetEntitlementsBy(ctx, entitlement.UserId)
	if err != nil {
		tx.RollbackTransaction()
		return err
	}

	err = tx.UpsertEntitlement(ctx, entitlement)
	if err != nil {
		tx.RollbackTransaction()
		return err
	}

	err = tx.CreateEntitlementHistory(ctx, &model.UserEntitlementHistory{
		UserEntitlementId: entitlement.ID,
		OldValue:          getValueForEntitlementBy(entitlement.Name, entitlements),
		NewValue:          entitlement.Value,
		Reason:            reason,
	})
	if err != nil {
		tx.RollbackTransaction()
		return err
	}

	tx.CommitTransaction()
	return nil
}

func (s *service) IsEntitlementEnabled(ctx context.Context, portfolioId uuid.UUID, name string) (bool, error) {
	entitlements, err := s.repository.GetEntitlementsBy(ctx, portfolioId)
	if err != nil {
		return false, err
	}

	return getValueForEntitlementBy(name, entitlements), nil
}

func getValueForEntitlementBy(name string, entitlements []model.UserEntitlement) bool {
	for _, value := range entitlements {
		if name == value.Name {
			return value.Value
		}
	}

	return false
}

func NewService(repository respository.Repository) Service {
	return &service{
		repository: repository,
	}
}
