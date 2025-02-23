package main

import (
	"easystore/configs/env"
	"easystore/db"
	"easystore/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	env.Load()
	db.Connect()
}

// @title Superstore API Docs
// @version 1.0
// @description API documentation for Superstore backend APIs
// @host localhost:8080
// @BasePath /api/v1
func main() {
	if os.Getenv("ENV") == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	routes.Intiliaze(r)
	r.Run(":8080")
}
