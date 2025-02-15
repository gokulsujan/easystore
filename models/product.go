package models

import (

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	OutletId         uint            `json:"outlet_id" gorm:"not null"`
	Outlet           Outlet          `gorm:"foreignKey:OutletId"`
	Title            string          `json:"title" gorm:"not null"`
	Description      string          `json:"description" gorm:"not null"`
	CategoryId       uint            `json:"category_id" gorm:"not null"`
	Category         ProductCategory `json:"foreignKey:CategoryId"`
	Status           string          `json:"status" gorm:"not null"`
}
