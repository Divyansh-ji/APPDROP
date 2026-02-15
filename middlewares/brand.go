package middlewares

import (
	"APPDROP/db"
	"APPDROP/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyBrandID = "brand_id"
	ContextKeyBrand   = "brand"
)

func BrandResolver() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := resolveDomain(c)
		if domain == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": gin.H{"code": "VALIDATION_ERROR", "message": "Brand domain required (subdomain or X-Brand-Domain header)"},
			})
			return
		}
		var brand models.Brand
		if err := db.DB.Where("domain = ?", domain).First(&brand).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Brand not found for domain: " + domain},
			})
			return
		}
		c.Set(ContextKeyBrandID, brand.ID)
		c.Set(ContextKeyBrand, &brand)
		c.Next()
	}
}

func resolveDomain(c *gin.Context) string {
	if h := c.GetHeader("X-Brand-Domain"); h != "" {
		return strings.TrimSpace(strings.ToLower(h))
	}

	host := c.Request.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}
	parts := strings.SplitN(host, ".", 2)
	if len(parts) < 2 {
		return ""
	}
	sub := strings.TrimSpace(strings.ToLower(parts[0]))
	if sub == "" || sub == "www" {
		return ""
	}
	return sub
}
