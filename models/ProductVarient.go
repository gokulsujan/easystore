package models

import "gorm.io/gorm"

type ProductVarient struct {
	gorm.Model
	ProductId    uint    `json:"product_id" gorm:"not null"`
	Product      Product `gorm:"foreignKey:ProductId"`
	Name         string  `json:"name" gorm:"not null"`
	SellingPrice float64 `json:"selling_price" gorm:"not null;type:decimal(10,2)"`
	Mrp          float64 `json:"mrp" gorm:"not null;type:decimal(10,2)"`
}
