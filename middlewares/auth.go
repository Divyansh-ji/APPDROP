package middlewares

import (
	"APPDROP/auth"
	"APPDROP/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		brandVal, exists := c.Get(ContextKeyBrand)
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Brand not found for this domain"},
			})
			return
		}
		brand, ok := brandVal.(*models.Brand)
		if !ok || brand == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": gin.H{"code": "NOT_FOUND", "message": "Brand not found"},
			})
			return
		}

		cookie, err := c.Cookie(auth.CookieName())
		if err != nil || strings.TrimSpace(cookie) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid session"},
			})
			return
		}

		claims, err := auth.ParseAndValidate(cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid or expired session"},
			})
			return
		}

		// claims.BrandID is already uuid.UUID, so no need to parse
		if claims.BrandID != brand.ID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": gin.H{"code": "FORBIDDEN", "message": "Brand in session does not match this domain"},
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
