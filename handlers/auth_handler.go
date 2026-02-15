package handlers

import (
	"APPDROP/auth"
	"APPDROP/db"
	"APPDROP/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	brand, ok := getBrandFromContext(c)
	if !ok || brand == nil {

		if brandID, ok := getBrandID(c); ok {
			var b models.Brand
			if err := db.DB.First(&b, "id = ?", brandID).Error; err == nil {
				brand = &b
				ok = true
			}
		}
	}
	if !ok || brand == nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Brand not found for this domain")
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	expectPass := os.Getenv("LOGIN_PASSWORD")
	if expectPass == "" {
		expectPass = "admin123"
	}
	if req.Password != expectPass {
		RespondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials")
		return
	}

	token, err := auth.CreateToken(brand.ID, req.Email, auth.DefaultTokenDuration)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create session")
		return
	}

	auth.SetSessionCookie(c.Writer, token, c.Request.Host, false)
	c.JSON(http.StatusOK, gin.H{"message": "ok", "brand_id": brand.ID.String()})
}

func Logout(c *gin.Context) {
	auth.ClearSessionCookie(c.Writer, c.Request.Host, false)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
