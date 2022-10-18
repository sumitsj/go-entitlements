package gorm

import (
	"fmt"
	"github.com/sumitsj/go-entitlements/constant"
	"github.com/sumitsj/go-entitlements/gorm/model"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) error {
	tx := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", constant.SchemaName))
	if tx.Error != nil {
		return tx.Error
	}
	return db.AutoMigrate(&model.UserEntitlement{}, &model.UserEntitlementHistory{})
}
