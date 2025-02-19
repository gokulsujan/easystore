package auth

import (
	"easystore/db"
	"easystore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var outlet models.Outlet
var outletEmployee models.OutletEmployee

func OutletMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		outlet_id := c.Param("outlet_id") // Extract tenant_id from URL param

		if outlet_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Tenant ID is required"})
			c.Abort()
			return
		}

		tx := db.DB.First(&outletEmployee, CurrentUserID(c))
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get user details"})
			return
		}

		tx = db.DB.First(&outlet, outlet_id)
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get outlet details"})
			return
		}

		if outlet.ID != outletEmployee.OutletId {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid outlet id"})
			c.Abort()
			return
		}

		c.Set("outlet_id", outlet_id)

		c.Next()
	}
}
