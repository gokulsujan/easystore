package routes

import (
	outlet_handler "easystore/handlers/outlet"
	_ "easystore/docs" // Import docs

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
	outletRoutes.POST("", outlet_handler.Create)
}
