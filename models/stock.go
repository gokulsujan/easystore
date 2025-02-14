package models

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	VarientId      uint           `json:"varient_id" gorm:"not null"`
	ProductVarient ProductVarient `gorm:"foreignKey:VarientId"`
	Quantity       int            `json:"quantity" gorm:"not null"`
}
