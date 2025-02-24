package routes

import (
	"easystore/auth"
	_ "easystore/docs"
	employeeHandler "easystore/handlers/employee"
	outletHandler "easystore/handlers/outlet"
	"easystore/handlers/product_category_handler"
	"easystore/handlers/product_varient_handler"
	product_handler "easystore/handlers/products"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


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

	productCategoryRoutes := outletRoutes.Group("/:outlet_id/product-category")
	productCategoryRoutes.Use(auth.OutletMiddleware())
	productCategoryRoutes.POST("", product_category_handler.Create)
	productCategoryRoutes.GET("/:category_id", product_category_handler.GetProductCategoryDetail)
	productCategoryRoutes.GET("", product_category_handler.GetProductCategories)
	productCategoryRoutes.PUT("/:category_id", product_category_handler.Update)

	productVarientRoutes := productRoutes.Group("/:product_id/product-varient")
	productVarientRoutes.POST("", product_varient_handler.Create)
	productVarientRoutes.PUT("/:varient_id", product_varient_handler.Update)
	productVarientRoutes.GET("", product_varient_handler.GetProductVarients)
	productVarientRoutes.GET("/:varient_id", product_varient_handler.GetProductVarient)

}
