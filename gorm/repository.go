package entitlements

import (
	"context"
	"github.com/google/uuid"
	"github.com/sumitsj/go-entitlements/gorm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetEntitlementsBy(ctx context.Context, userId uuid.UUID) ([]model.UserEntitlement, error)
	UpsertEntitlement(ctx context.Context, entitlement *model.UserEntitlement) error
	CreateEntitlementHistory(ctx context.Context, entitlementHistory *model.UserEntitlementHistory) error

	WithTransaction() Repository
	CommitTransaction()
	RollbackTransaction()
}

func (r *repository) GetEntitlementsBy(ctx context.Context, userId uuid.UUID) ([]model.UserEntitlement, error) {
	var entitlements []model.UserEntitlement
	tx := r.db.WithContext(ctx).Find(&entitlements, model.UserEntitlement{UserId: userId})
	return entitlements, tx.Error
}

func (r *repository) UpsertEntitlement(ctx context.Context, entitlement *model.UserEntitlement) error {
	tx := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "name"}},
		UpdateAll: true,
	}).Create(entitlement)

	return tx.Error
}

func (r *repository) CreateEntitlementHistory(ctx context.Context, entitlementHistory *model.UserEntitlementHistory) error {
	tx := r.db.WithContext(ctx).Create(entitlementHistory)

	return tx.Error
}

func (r *repository) WithTransaction() Repository {
	return &repository{db: r.db.Begin()}
}

func (r *repository) CommitTransaction() {
	r.db.Commit()
}

func (r *repository) RollbackTransaction() {
	r.db.Rollback()
}

//goland:noinspection GoUnusedExportedFunction
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
