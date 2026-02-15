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

	routes.RegisterRoutes(r)

	log.Println("Server running on port 8090")
	r.Run(":8090")
}
