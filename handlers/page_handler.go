package handlers

import (
	"APPDROP/db"
	"APPDROP/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePages(c *gin.Context) {
	var page models.Page

	if err := c.ShouldBindJSON(&page); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if page.Name == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "page name is required")
		return
	}

	if page.Route == "" {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "page route is required")
		return
	}

	var existing models.Page
	if err := db.DB.Where("route = ?", page.Route).First(&existing).Error; err == nil {
		RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "Page route already exists")
		return
	}
	if page.IsHome {
		var homePage models.Page
		if err := db.DB.Where("is_home = true").First(&homePage).Error; err == nil {
			RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "home page already exists")
			return
		}
	}
	if err := db.DB.Create(&page).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create a page")
		return
	}
	c.JSON(http.StatusCreated, page)

}

func GetPages(c *gin.Context) {
	var pages []models.Page

	if err := db.DB.Find(&pages).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch pages")
		return
	}

	c.JSON(http.StatusOK, pages)
}

func GetPageByID(c *gin.Context) {
	idParam := c.Param("id")

	pageID, err := uuid.Parse(idParam)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page ID")
		return
	}

	var page models.Page
	if err := db.DB.Preload("Widgets").First(&page, "id = ?", pageID).Error; err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Page not found")
		return
	}

	c.JSON(http.StatusOK, page)
}

func UpdatePage(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page ID")
		return
	}

	var page models.Page
	if err := db.DB.First(&page, "id = ?", pageID).Error; err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Page not found")
		return
	}

	var input struct {
		Name   *string `json:"name"`
		Route  *string `json:"route"`
		IsHome *bool   `json:"is_home"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if input.Name != nil {
		if *input.Name == "" {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "page name is required")
			return
		}
		page.Name = *input.Name
	}
	if input.Route != nil {
		if *input.Route == "" {
			RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "page route is required")
			return
		}
		var existing models.Page
		if err := db.DB.Where("route = ? AND id != ?", *input.Route, pageID).First(&existing).Error; err == nil {
			RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "Page route already exists")
			return
		}
		page.Route = *input.Route
	}
	if input.IsHome != nil && *input.IsHome {
		var homePage models.Page
		if err := db.DB.Where("is_home = true AND id != ?", pageID).First(&homePage).Error; err == nil {
			RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "home page already exists")
			return
		}
		page.IsHome = true
	} else if input.IsHome != nil {
		page.IsHome = false
	}

	if err := db.DB.Save(&page).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to update page")
		return
	}
	c.JSON(http.StatusOK, page)
}
