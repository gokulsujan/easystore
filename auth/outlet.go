package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OutletMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		outlet_id := c.Param("outlet_id") // Extract tenant_id from URL param

		if outlet_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Tenant ID is required"})
			c.Abort()
			return
		}

		c.Set("outlet_id", outlet_id)

		c.Next()
	}
}
