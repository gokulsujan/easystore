package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Phone    string `json:"phone" gorm:"not null;size:10;unique"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	Role     string `json:"role" gorm:"not null"`
	Status   string `json:"status" gorm:"not null"`
}

// HashPassword hashes the password of an employee
func (e *Employee) HashPassword() error {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(bcryptPassword)
	return nil
}
