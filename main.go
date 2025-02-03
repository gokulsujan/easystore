package main

import (
	"easystore/configs/env"
	"easystore/db"

	"github.com/gin-gonic/gin"
)

func init() {
	env.Load()
	db.Connect()
}

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Pong",
		})
	})
	r.Run(":8080")
}
