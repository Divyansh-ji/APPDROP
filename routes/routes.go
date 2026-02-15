package routes

import (
	"APPDROP/handlers"
	"APPDROP/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Public (no brand required)
	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.POST("/brands", handlers.CreateBrand)

	// Brand-scoped (BrandResolver: brand from subdomain or X-Brand-Domain)
	brandGroup := r.Group("/")
	brandGroup.Use(middlewares.BrandResolver())
	{
		// Public within brand: login (sets cookie for this brand's domain)
		brandGroup.POST("/login", handlers.Login)
		brandGroup.POST("/logout", handlers.Logout)

		// Protected: require valid JWT and brand match
		protected := brandGroup.Group("/")
		protected.Use(middlewares.RequireAuth())
		{
			protected.POST("/pages", handlers.CreatePages)
			protected.GET("/pages", handlers.GetPages)
			protected.GET("/pages/:id", handlers.GetPageByID)
			protected.PUT("/pages/:id", handlers.UpdatePage)
			protected.DELETE("/pages/:id", handlers.DeletePage)
			protected.POST("/pages/:id/widgets", handlers.AddWidget)
			protected.PUT("/widgets/:id", handlers.UpdateWidget)
			protected.DELETE("/widgets/:id", handlers.DeleteWidget)
			protected.POST("/pages/:id/widgets/reorder", handlers.ReorderWidgets)
			protected.GET("/brands/me", handlers.GetBrandMe)
			protected.GET("/brands/:id", handlers.GetBrandByID)
		}
	}
}
