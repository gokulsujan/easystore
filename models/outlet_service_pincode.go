package models

import "gorm.io/gorm"

type OutletServicePincode struct {
	gorm.Model
	OutletId uint   `json:"outlet_id" gorm:"not null"`
	Outlet   Outlet `gorm:"foreignKey:OutletId"`
	Pincode  string `json:"pincode", gorm:"unique"`
}
