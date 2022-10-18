package model

import (
	"github.com/sumitsj/go-entitlements/constant"
	"gorm.io/gorm"
)

type UserEntitlementHistory struct {
	gorm.Model
	UserEntitlementId uint
	OldValue          bool
	NewValue          bool
	Reason            string
}

func (UserEntitlementHistory) TableName() string {
	return constant.SchemaName + ".user_entitlements_history"
}
