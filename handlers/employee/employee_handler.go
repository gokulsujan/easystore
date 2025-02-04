package employee_handler

import (
	"easystore/db"
	"easystore/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var employee models.Employee

// CreateEmployee is a http request handler which creates a new employee
func CreateEmployee(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":err.Error()})
		return
	}

	if !validEmployeeFields(employee, c) {
		return
	}

	employee.Status = "active"
	err = employee.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed", "message":"Failed to hash password"})
		return
	}

	// Save employee to database
	tx := db.DB.Create(&employee)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed", "message":"Failed to create employee"})
		return
	}
	c.JSON(200, gin.H{"status":"success","message": "Create employee", "result": employee})
}

// Private methods

var validEmployeeFields = func(employee models.Employee, c *gin.Context) bool {
	if employee.Name == "" || employee.Phone == "" || employee.Email == "" || employee.Password == "" || employee.Role == "" || employee.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"All fields are required"})
		return false
	}

	if len(employee.Phone) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Phone number must be 10 digits"})
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(employee.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Invalid email address"})
		return false
	}

	tx := db.DB.Where("email = ?", employee.Email).First(&employee)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Email already exists"})
		return false
	}

	tx = db.DB.Where("phone = ?", employee.Phone).First(&employee)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Phone number already exists"})
		return false
	}
	return true
}
