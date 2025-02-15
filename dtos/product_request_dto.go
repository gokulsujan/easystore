package dtos

import "time"

type Product struct {
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	CategoryId       uint             `json:"category_id"`
	ManufacturedDate time.Time        `json:"manufactured_date"`
	ExpiryDate       time.Time        `json:"expiry_date"`
	Status           string           `json:"status"`
	Varients         []ProductVarient `josn:"varients"`
}

type ProductVarient struct {
	ProductId    uint    `json:"product_id"`
	Name         string  `json:"name"`
	SellingPrice float64 `json:"selling_price"`
	Mrp          float64 `json:"mrp"`
}
