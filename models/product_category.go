package models

import "gorm.io/gorm"

type ProductCategory struct {
	gorm.Model
	OutletId    uint   `json:"outlet_id"`
	Outlet      Outlet `gorm:"foreignKey:OutletId"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
}
