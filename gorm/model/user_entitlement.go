package model

import (
	"github.com/google/uuid"
	"github.com/sumitsj/go-entitlements/constant"
	"gorm.io/gorm"
)

type UserEntitlement struct {
	gorm.Model
	UserId uuid.UUID `gorm:"index:,unique,composite:user_id_and_name"`
	Name   string    `gorm:"index:,unique,composite:user_id_and_name"`
	Value  bool

	UserEntitlementHistories []UserEntitlementHistory
}

func (UserEntitlement) TableName() string {
	return constant.SchemaName + ".user_entitlements"
}
