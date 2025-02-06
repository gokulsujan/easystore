package employee_handler

import (
	"easystore/db"
	"easystore/dtos"
	handler_helper "easystore/handlers/helpers"
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
// @Param Authorization header string true "Bearer Token"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /employee [post]
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
// @Param Authorization header string true "Bearer Token"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /employee/{id} [put]
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

// @Summary      Login an employee
// @Description  Logs in an employee and returns an access token
// @Tags         Employee
// @Accept       json
// @Produce      json
// @Param        employee  body  dtos.EmployeeLogin  true  "Employee Login Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      401  {object}  dtos.ErrorResponse
// @Failure      404  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Router       /employee/login [post]
func Login(c *gin.Context) {
	var employeeLogin dtos.EmployeeLogin
	err := c.ShouldBindBodyWithJSON(&employeeLogin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	if employeeLogin.Email == "" || employeeLogin.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Email and password are required"})
		return
	}

	// Find employee by email
	tx := db.DB.Where("email = ?", employeeLogin.Email).First(&employee)
	if tx.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Employee not found"})
		return
	}

	// Verify password
	if !employee.VerifyPassword(employeeLogin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid email or password"})
		return
	}

	employee.OmitPassword()
	token, err := handler_helper.GenerateEmployeeLoginJwt(&employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to generate token", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Login successful", "result": gin.H{"accessToken": token}})
}

// @Summary      Get an employee
// @Description  Gets an employee by ID
// @Param  id path string true "Employee ID"
// @Param Authorization header string true "Bearer Token"
// @Tags         Employee
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      404  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /employee/{id} [get]
func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	tx := db.DB.Where("id = ?", id).First(&employee)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Employee not found"})
		return
	} else if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to get employee"})
		return
	}
	employee.OmitPassword()
	c.JSON(200, gin.H{"status": "success", "message": "Get employee", "result": gin.H{"employee": employee}})
}

// @Summary      Get all employees
// @Description  Gets all employees
// @Param Authorization header string true "Bearer Token"
// @Tags         Employee
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /employee [get]
func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	tx := db.DB.Omit("password").Find(&employees)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to get employees"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Get employees", "result": gin.H{"employees": employees}})
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
