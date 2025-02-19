package routes

import (
	"easystore/auth"
	_ "easystore/docs"
	employeeHandler "easystore/handlers/employee"
	outletHandler "easystore/handlers/outlet"
	product_handler "easystore/handlers/products"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Superstore API Docs
// @version 1.0
// @description Api documentation for Superstore backend apis
// @host localhost:8080
// @BasePath /api/v1

func Intiliaze(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api/v1")
	api.POST("/employee/login", employeeHandler.Login)

	outletRoutes := api.Group("/outlet")
	outletRoutes.Use(auth.JWTMiddleware())
	outletRoutes.POST("", outletHandler.Create)
	outletRoutes.PUT("/:outlet_id", outletHandler.Update)
	outletRoutes.GET("", outletHandler.GetOutlets)
	outletRoutes.GET("/:outlet_id", outletHandler.GetOutlet)
	outletRoutes.POST("/:outlet_id/assign-pincodes", outletHandler.AssignOutletServicePincode)

	employeeRoutes := api.Group("/employee")
	employeeRoutes.Use(auth.JWTMiddleware())
	employeeRoutes.POST("", employeeHandler.Create)
	employeeRoutes.PUT("/:employee_id", employeeHandler.Update)
	employeeRoutes.GET("", employeeHandler.GetEmployees)
	employeeRoutes.GET("/:employee_id", employeeHandler.GetEmployee)
	employeeRoutes.POST("/:employee_id/outlet", employeeHandler.CreateOutlet)

	productRoutes := outletRoutes.Group(("/:outlet_id/product"))
	productRoutes.Use(auth.OutletMiddleware())
	productRoutes.GET("/:product_id", product_handler.GetProductDetails)
	productRoutes.POST("", product_handler.Create)
	productRoutes.PUT("/:product_id", product_handler.Update)
}
