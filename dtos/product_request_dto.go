package dtos

type Product struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	CategoryId  uint             `json:"category_id"`
	Status      string           `json:"status"`
	Varients    []ProductVarient `josn:"varients"`
}

type ProductVarient struct {
	ProductId    uint    `json:"product_id"`
	Name         string  `json:"name"`
	SellingPrice float64 `json:"selling_price"`
	Mrp          float64 `json:"mrp"`
}

type ProductCategory struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
