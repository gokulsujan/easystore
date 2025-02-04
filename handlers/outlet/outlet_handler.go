package outlet_handler

import (
	"easystore/db"
	"easystore/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var outlet models.Outlet

func Create(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&outlet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed", "message":err.Error()})
		return
	}

	if !validOutletFields(outlet, c) {
		return
	}

	// Save outlet to database
	tx := db.DB.Create(&outlet)

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed", "message":"Failed to create outlet"})
		return
	}
	c.JSON(200, gin.H{"status":"success","message": "Create outlet", "result": outlet})
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
	return true
}
