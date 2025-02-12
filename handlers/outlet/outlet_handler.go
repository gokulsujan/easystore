package outlet_handler

import (
	"easystore/db"
	"easystore/dtos"
	handler_helper "easystore/handlers/helpers"
	"easystore/models"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var outlet models.Outlet

// @Summary      Create an outlet
// @Description  Creates a new outlet and returns the created outlet object
// @Param Authorization header string true "Bearer Token"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Outlet  true  "Outlet Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet [post]
func Create(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&outlet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	if !validOutletFields(outlet, c) {
		return
	}

	// Generate unique identifier for outlet
	outlet.Identifier = handler_helper.GenerateUUID()

	// Save outlet to database
	tx := db.DB.Create(&outlet)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to create outlet"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Create outlet", "result": outlet})
}

// @Summary      Update an outlet
// @Description  Updates an existing outlet and returns the updated outlet object
// @Param Authorization header string true "Bearer Token"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Outlet  true  "Outlet Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet [put]
func Update(c *gin.Context) {
	id := c.Param("id")
	err := c.ShouldBindBodyWithJSON(&outlet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	if outlet.Name == "" && outlet.Description == "" && outlet.Location == "" && outlet.Phone == "" && outlet.Email == "" && outlet.Website == "" && outlet.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Atleast one field is required"})
		return
	}

	// Save outlet to database
	tx := db.DB.Model(models.Outlet{}).Where("id = ?", id).Updates(&outlet)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to update outlet"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Update outlet", "result": outlet})
}

// @Summary      Assign manager to outlet
// @Description  Assigns a manager to an outlet and returns the updated outlet object
// @Param Authorization header string true "Bearer Token"
// @Param  outlet_id path string true "Outlet ID"
// @Param  manager_id query string true "Manager ID"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/assign-manager [put]
func AssignManager(c *gin.Context) {
	outletId := c.Param("id")
	managerId := c.Query("manager_id")

	// Check if outlet ID and manager ID are provided
	if outletId == "" || managerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Outlet ID and Manager ID are required"})
		return
	}

	// Check if outlet exists
	outletTx := db.DB.Where("id = ?", outletId).First(&outlet)
	if outletTx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Outlet not found with the given ID"})
		return
	}

	// Check if manager exists
	var manager models.Employee
	managerTx := db.DB.Where("id = ?", managerId).Omit("password").First(&manager)
	if managerTx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Manager not found with the given ID"})
		return
	}

	// Assign manager to outlet
	outlet.ManagerId = manager.ID
	tx := db.DB.Model(&outlet).Updates(&outlet)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to assign manager to outlet"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "Assign manager to outlet", "result": gin.H{"outlet": outlet, "manager": manager}})
}

// @Summary      Get all outlets
// @Description  Returns a list of all outlets
// @Param Authorization header string true "Bearer Token"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet [get]
func GetOutlets(c *gin.Context) {
	var outlets []models.Outlet
	tx := db.DB.Preload("Manager", func(db *gorm.DB) *gorm.DB {
							return db.Omit("password")
						}).Find(&outlets)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to fetch outlets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Get outlets", "result": gin.H{"outlets": outlets}})
}

// @Summary      Get an outlet
// @Description  Gets an outlet by ID
// @Param  id path string true "Outlet ID"
// @Param Authorization header string true "Bearer Token"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{id} [get]
func GetOutlet(c *gin.Context) {
	id := c.Param("id")
	tx := db.DB.Where("id = ?", id).Preload("Manager", func(db *gorm.DB) *gorm.DB {
												return db.Omit("password")
											}).First(&outlet)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to fetch outlet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Get outlet", "result": outlet})
}

// @Summary      Assign Pincodes
// @Description  Assign service area pincodes of an outlet
// @Param  id path string true "Outlet ID"
// @Param Authorization header string true "Bearer Token"
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Param        pincodes  body  dtos.OutletPincodes  true  "Service Pincodes"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{id}/assign-pincodes [get]
func AssignOutletServicePincode(c *gin.Context) {
	var pincodes dtos.OutletPincodes
	c.ShouldBindBodyWithJSON(&pincodes)

	outletId := c.Param("id")

	// Check if outlet ID pincodes are provided
	if len(pincodes.Pincodes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Pincodes are required"})
		return
	}

	// Find outlet
	tx := db.DB.First(&outlet)
	if tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Outlet not found with the given ID"})
		return
	}

	outletIdNum, err := strconv.Atoi(outletId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to convert outlet id to number", "result": gin.H{"error":"error"}})
		return
	}

	var failedPincodes []string
	var successPincodes []string
	for _, pincode := range pincodes.Pincodes {
		var outletServicePincode models.OutletServicePincode
		outletServicePincode.OutletId =uint(outletIdNum)
		outletServicePincode.Pincode = pincode
		tx := db.DB.Create(outletServicePincode)
		if tx.Error == nil {
			successPincodes = append(successPincodes, pincode)
		} else {
			failedPincodes = append(failedPincodes, pincode)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message":"Pincodes assigned to the outlet", "result": gin.H{"success":successPincodes, "failed": failedPincodes}})
}

// Private methods

var validOutletFields = func(outlet models.Outlet, c *gin.Context) bool {
	if outlet.Name == "" || outlet.Description == "" || outlet.Location == "" || outlet.Phone == "" || outlet.Email == "" || outlet.Website == "" || outlet.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "All fields are required"})
		return false
	}

	if len(outlet.Phone) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Phone number must be 10 digits"})
		return false
	}

	if outlet.Status != "active" && outlet.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Status must be active or inactive"})
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(outlet.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid email address"})
		return false
	}

	tx := db.DB.Where("email = ?", outlet.Email).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Email already exists"})
		return false
	}

	tx = db.DB.Where("phone = ?", outlet.Phone).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Phone number already exists"})
		return false
	}

	tx = db.DB.Where(("website = ?"), outlet.Website).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Website already exists"})
		return false
	}
	return true
}
