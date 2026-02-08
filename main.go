package main

import (
	"APPDROP/db"
	"APPDROP/middlewares"
	"APPDROP/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	r.Use(middlewares.RequestLogger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	routes.RegisterRoutes(r)

	log.Println(" Server running on port 8082")
	r.Run(":8082")
}
