package routes

import (
	_ "easystore/docs" // Import docs
	employeeHandler "easystore/handlers/employee"
	outletHandler "easystore/handlers/outlet"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Superstore API Docs
// @version 1.0
// @description Api documentation for Superstore backend apis
// @host localhost:8080
// @BasePath /api/v1

func Intiliaze(r *gin.Engine) {
	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	outletRoutes := api.Group("/outlet")
	outletRoutes.POST("", outletHandler.Create)
	outletRoutes.PUT("/:id", outletHandler.Update)

	employeeRoutes := api.Group("/employee")
	employeeRoutes.POST("", employeeHandler.Create)
	employeeRoutes.PUT("/:id", employeeHandler.Update)
}
