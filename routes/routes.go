package routes

import (
	"APPDROP/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Pages
	r.POST("/pages", handlers.CreatePages)
	r.GET("/pages", handlers.GetPages)
	r.GET("/pages/:id", handlers.GetPageByID)
	r.DELETE("/pages/:id", handlers.DeletePage)

	// Widgets
	r.POST("/pages/:id/widgets", handlers.AddWidget)
	r.PUT("/widgets/:id", handlers.UpdateWidget)
	r.DELETE("/widgets/:id", handlers.DeleteWidget)
	r.POST("/pages/:id/widgets/reorder", handlers.ReorderWidgets)
}
