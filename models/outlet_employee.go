package models

type OutletEmployee struct {
	OutletId uint `json:"outlet_id" gorm:"not null"`
	Outlet Outlet `gorm:"foreignKey:OutletId"`
    EmployeeId uint `json:"employee_id" gorm:"not null"`
    Employee Employee `gorm:"foreignKey:EmployeeId"`
    Role string `json:"role" gorm:"not null"`
}
