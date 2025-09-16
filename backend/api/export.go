package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// ExportModels serves the entire model catalog as a downloadable JSON payload.
// It expects no parameters on the incoming request, queries the database for
// models with nested versions/images, and writes the JSON export with
// Content-Disposition headers. The handler is read-only with respect to the
// database but streams the response body to the client.
func ExportModels(c *gin.Context) {
	var modelsList []models.Model
	database.DB.Preload("Versions").Preload("Versions.Images").Find(&modelsList)

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=\"model_export.json\"")
	c.JSON(http.StatusOK, modelsList)
}
