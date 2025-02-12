package models

import "gorm.io/gorm"

type Outlet struct {
	gorm.Model
	Identifier  string   `json:"identifier" gorm:"unique;not null"`
	Name        string   `json:"name" gorm:"unique"`
	Description string   `json:"description" gorm:"not null"`
	ManagerId   uint     `json:"manager_id" gorm:"not null"`
	Manager     Employee `gorm:"foreignKey:ManagerId"`
	Location    string   `json:"location" gorm:"not null"`
	Phone       string   `json:"phone" gorm:"not null;size:10;unique"`
	Email       string   `json:"email" gorm:"not null"`
	Website     string   `json:"website" gorm:"not null"`
	Status      string   `json:"status" gorm:"not null"`
}
