package handlers

import (
	"APPDROP/db"
	"APPDROP/middlewares"
	"APPDROP/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateBrandRequest struct {
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	OfficeAddress string `json:"office_address"`
	Logo          string `json:"logo"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}

func CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if req.Name == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "brand name is required")
		return
	}

	if req.Domain == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "brand domain is required")
		return
	}

	if req.Email == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "brand email is required")
		return
	}

	if req.Password == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "brand password is required")
		return
	}

	var existingBrand models.Brand

	if err := db.DB.Where("domain = ?", req.Domain).First(&existingBrand).Error; err == nil {
		RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "brand domain already exists")
		return
	}

	if err := db.DB.Where("email = ?", req.Email).First(&existingBrand).Error; err == nil {
		RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "brand email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to hash password")
		return
	}

	brand := models.Brand{
		Name:          req.Name,
		Domain:        req.Domain,
		OfficeAddress: req.OfficeAddress,
		Logo:          req.Logo,
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
	}

	if err := db.DB.Create(&brand).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create a brand")
		return
	}
	c.JSON(http.StatusCreated, brand)
}
func GetBrandByID(c *gin.Context) {
	brandVal, exists := c.Get(middlewares.ContextKeyBrand)
	if !exists {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Brand not found for this domain")
		return
	}
	ctxBrand, ok := brandVal.(*models.Brand)
	if !ok || ctxBrand == nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Brand not found")
		return
	}
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid brand ID")
		return
	}
	if id != ctxBrand.ID {
		RespondError(c, http.StatusForbidden, "FORBIDDEN", "Cannot access another brand")
		return
	}
	c.JSON(http.StatusOK, ctxBrand)
}
func GetBrandMe(c *gin.Context) {
	brandVal, exists := c.Get(middlewares.ContextKeyBrand)
	if !exists {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Brand not found for this domain")
		return
	}
	brand, ok := brandVal.(*models.Brand)
	if !ok || brand == nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Brand not found")
		return
	}
	c.JSON(http.StatusOK, brand)
}
