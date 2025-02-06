package outlet_handler

import (
	"easystore/db"
	"easystore/models"
	handler_helper "easystore/handlers/helpers"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var outlet models.Outlet

// @Summary      Create an outlet
// @Description  Creates a new outlet and returns the created outlet object
// @Tags         Outlet
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Outlet  true  "Outlet Details"
// @Param Authorization header string true "Bearer Token"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /api/v1/outlet [post]
func Create(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&outlet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed", "message":"Failed to create outlet"})
		return
	}
	c.JSON(200, gin.H{"status":"success","message": "Create outlet", "result": outlet})
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
// @Router       /api/v1/outlet [put]
func Update(c *gin.Context) {
	id := c.Param("id")
	err := c.ShouldBindBodyWithJSON(&outlet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":err.Error()})
		return
	}

	if outlet.Name == "" && outlet.Description == "" && outlet.Location == "" && outlet.Phone == "" && outlet.Email == "" && outlet.Website == "" && outlet.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Atleast one field is required"})
		return
	}

	// Save outlet to database
	tx := db.DB.Model(models.Outlet{}).Where("id = ?", id).Updates(&outlet)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed", "message":"Failed to update outlet"})
		return
	}
	c.JSON(200, gin.H{"status":"success","message": "Update outlet", "result": outlet})
}

// Private methods

var validOutletFields = func(outlet models.Outlet, c *gin.Context) bool {
	if outlet.Name == "" || outlet.Description == "" || outlet.Location == "" || outlet.Phone == "" || outlet.Email == "" || outlet.Website == "" || outlet.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"All fields are required"})
		return false
	}

	if len(outlet.Phone) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Phone number must be 10 digits"})
		return false
	}

	if outlet.Status != "active" && outlet.Status != "inactive" {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Status must be active or inactive"})
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(outlet.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Invalid email address"})
		return false
	}

	tx := db.DB.Where("email = ?", outlet.Email).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Email already exists"})
		return false
	}

	tx = db.DB.Where("phone = ?", outlet.Phone).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Phone number already exists"})
		return false
	}

	tx = db.DB.Where(("website = ?"), outlet.Website).First(&outlet)
	if tx.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":"Website already exists"})
		return false
	}
	return true
}
