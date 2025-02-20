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

func main() {
	if os.Getenv("ENV") == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	routes.Intiliaze(r)
	r.Run(":8080")
}
