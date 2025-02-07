package dtos

type Outlet struct {
	Name        string `json:"name" example:"Superstore Attingal"`
	Description string `json:"description" example:"Superstore Attingal is a supermarket located in Attingal, Kerala"`
	Location    string `json:"location" example:"Attingal, Kerala"`
	Phone       string `json:"phone" example:"9876543210"`
	Email       string `json:"email" example:"attingal@superstore.com"`
	Website     string `json:"website" example:"attingal.superstore.com"`
	Status      string `json:"status" example:"active"`
}

type OutletPincodes struct {
	Pincodes []string `json:"pincodes" example:"["695606", 695101", "695103"]"`
}
