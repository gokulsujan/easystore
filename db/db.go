package db

import (
	"easystore/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := os.Getenv("DB_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB.AutoMigrate(&models.Outlet{})
	DB.AutoMigrate(&models.Employee{})
	DB.AutoMigrate(&models.OutletEmployee{})
	DB.AutoMigrate(&models.OutletServicePincode{})
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.ProductCategory{})
	DB.AutoMigrate(&models.ProductVarient{})
	DB.AutoMigrate(&models.Stock{})
}
