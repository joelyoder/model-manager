package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// ExportModels returns a JSON representation of all models in the database
// including their versions and associated images.
func ExportModels(c *gin.Context) {
	var modelsList []models.Model
	database.DB.Preload("Versions").Preload("Versions.Images").Find(&modelsList)

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=\"model_export.json\"")
	c.JSON(http.StatusOK, modelsList)
}
