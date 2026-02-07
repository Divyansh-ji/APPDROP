package main

import (
	"APPDROP/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("🚀 Server running on port 8080")
	r.Run(":8082")
}
