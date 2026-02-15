package handlers

import (
	"APPDROP/middlewares"
	"APPDROP/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getBrandID returns the current brand ID from context (set by BrandResolver). ok is false if missing.
func getBrandID(c *gin.Context) (uuid.UUID, bool) {
	v, exists := c.Get(middlewares.ContextKeyBrandID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

// getBrandFromContext returns the current brand from context (set by BrandResolver). ok is false if missing.
func getBrandFromContext(c *gin.Context) (*models.Brand, bool) {
	v, exists := c.Get(middlewares.ContextKeyBrand)
	if !exists {
		return nil, false
	}
	brand, ok := v.(*models.Brand)
	return brand, ok
}
