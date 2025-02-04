package main

import (
	"easystore/configs/env"
	"easystore/db"
	"easystore/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	env.Load()
	db.Connect()
}

func main() {

	r := gin.Default()

	routes.Intiliaze(r)
	r.Run(":8080")
}
