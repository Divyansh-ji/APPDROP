package handlers

import (
	"APPDROP/db"
	"APPDROP/models"
	"net/http"
	"strconv"

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
	pageParam := c.Query("page")
	limitParam := c.Query("limit")

	if pageParam == "" && limitParam == "" {
		var pages []models.Page
		if err := db.DB.Find(&pages).Error; err != nil {
			RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch pages")
			return
		}
		c.JSON(http.StatusOK, pages)
		return
	}

	page := 1
	limit := 10
	if p := pageParam; p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := limitParam; l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100
			}
		}
	}
	offset := (page - 1) * limit

	var total int64
	if err := db.DB.Model(&models.Page{}).Count(&total).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to count pages")
		return
	}
	var pages []models.Page
	if err := db.DB.Order("created_at ASC").Offset(offset).Limit(limit).Find(&pages).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch pages")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  pages,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func GetPageByID(c *gin.Context) {
	idParam := c.Param("id")

	pageID, err := uuid.Parse(idParam)
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page ID")
		return
	}

	widgetTypeFilter := c.Query("widget_type")

	if widgetTypeFilter != "" && !IsAllowedWidgetType(widgetTypeFilter) {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid widget_type filter")
		return
	}

	var page models.Page
	query := db.DB
	if widgetTypeFilter != "" {
		query = query.Preload("Widgets", "type = ?", widgetTypeFilter)
	} else {
		query = query.Preload("Widgets")
	}
	if err := query.First(&page, "id = ?", pageID).Error; err != nil {
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
