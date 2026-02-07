package routes

import (
	"APPDROP/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/pages", handlers.CreatePages)
	r.GET("/pages", handlers.GetPages)
	r.GET("/pages/:id", handlers.GetPageByID)
}
