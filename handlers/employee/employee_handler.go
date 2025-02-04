package employee_handler

import (
	"easystore/db"
	"easystore/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var employee models.Employee

// @Summary      Create an employee
// @Description  Creates a new employee and returns the created employee object
// @Tags         Employee
// @Accept       json
// @Produce      json
// @Param        employee  body  dtos.EmployeeCreate  true  "Employee Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Router       /api/v1/employee [post]

// CreateEmployee is a http request handler which creates a new employee
func Create(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	if !validEmployeeFields("create", employee, c) {
		return
	}

	employee.Status = "active"
	err = employee.HashPassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to hash password"})
		return
	}

	// Save employee to database
	tx := db.DB.Create(&employee)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to create employee"})
		return
	}
	employee.OmitPassword()
	c.JSON(200, gin.H{"status": "success", "message": "Create employee", "result": employee})
}

// @Summary      Update an employee
// @Description  Updates an existing employee and returns the updated employee object
// @Param  id path string true "Employee ID"
// @Tags         Employee
// @Accept       json
// @Produce      json
// @Param        employee  body  dtos.EmployeeUpdate  true  "Employee Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Router       /api/v1/employee/{id} [put]

// UpdateEmployee is a http request handler which updates an existing employee
func Update(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&employee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	id := c.Param("id")

	if !validEmployeeFields("update", employee, c) {
		return
	}
	if employee.Password != "" {
		err = employee.HashPassword()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to hash password"})
			return
		}
	}	
	// Save employee to database
	tx := db.DB.Model(&models.Employee{}).Where("id = ?", id).Updates(&employee)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to update employee"})
		return
	}

	employee.OmitPassword()
	c.JSON(200, gin.H{"status": "success", "message": "Update employee", "result": employee})
}

// Private methods

var validEmployeeFields = func(operation string, employee models.Employee, c *gin.Context) bool {
	if operation == "create" && (employee.Name == "" || employee.Phone == "" || employee.Email == "" || employee.Role == "" && employee.Password == "") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "All fields are required"})
		return false
	}

	if employee.Phone != "" && len(employee.Phone) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Phone number must be 10 digits"})
		return false
	}

	var tx *gorm.DB
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if employee.Email != "" {
		if !emailRegex.MatchString(employee.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid email address"})
			return false
		}

		if operation == "update" {
			tx = db.DB.Where("email = ?", employee.Email).Not("id = ?", employee.ID).First(&employee)
		} else {
			tx = db.DB.Where("email = ?", employee.Email).First(&employee)
		}

		if tx.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Email already exists"})
			return false
		}
	}

	if employee.Phone != "" {
		if operation == "update" {
			tx = db.DB.Where("phone = ?", employee.Phone).Not("id = ?", employee.ID).First(&employee)
		} else {
			tx = db.DB.Where("phone = ?", employee.Phone).First(&employee)
		}
		if tx.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Phone number already exists"})
			return false
		}
	}
	return true
}
