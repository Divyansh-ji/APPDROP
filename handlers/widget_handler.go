package handlers

import (
	"APPDROP/db"
	"APPDROP/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var allowedWidgetTypes = map[string]bool{
	"banner":       true,
	"product_grid": true,
	"text":         true,
	"image":        true,
	"spacer":       true,
}

func AddWidget(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page ID")
		return
	}

	var widget models.Widget
	if err := c.ShouldBindJSON(&widget); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if !allowedWidgetTypes[widget.Type] {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid widget type")
		return
	}

	widget.PageID = pageID

	if err := db.DB.Create(&widget).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create widget")
		return
	}
	c.JSON(http.StatusCreated, widget)
}
func UpdateWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid Widget ID")
		return
	}

	var widget models.Widget
	if err := db.DB.First(&widget, "id = ?", widgetID).Error; err != nil {
		RespondError(c, http.StatusNotFound, "NOT_FOUND", "Widget not found")
		return
	}
	if err := c.ShouldBindJSON(&widget); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	if !allowedWidgetTypes[widget.Type] {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid widget type")
		return
	}
	if err := db.DB.Save(&widget).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update the widget")
		return
	}
	c.JSON(http.StatusOK, widget)
}
func DeleteWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid widget ID")
		return
	}

	if err := db.DB.Delete(&models.Widget{}, "id = ?", widgetID).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete widget")
		return
	}

	c.Status(http.StatusNoContent)
}

type ReorderRequest struct {
	WidgetIDs []uuid.UUID `json:"widget_ids"`
}

func ReorderWidgets(c *gin.Context) {
	pageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page ID")
		return
	}
	var req ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	for index, widgetID := range req.WidgetIDs {
		db.DB.Model(&models.Widget{}).Where("id = ? AND page_id = ?", widgetID, pageID).Update("Position", index)
	}
	c.JSON(http.StatusOK, gin.H{"status": "reordered"})

}

func DeletePage(c *gin.Context) {
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

	if page.IsHome {
		RespondError(c, http.StatusConflict, "VALIDATION_ERROR", "Cannot delete home page")
		return
	}

	if err := db.DB.Delete(&page).Error; err != nil {
		RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete page")
		return
	}

	c.Status(http.StatusNoContent)
}
