package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetSettings loads every persisted application setting and returns them as a
// JSON map. The endpoint requires no parameters and only reads from the
// database.
func GetSettings(c *gin.Context) {
	settings, err := database.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// UpdateSetting upserts a single setting record based on the JSON payload in
// the request body. The body must provide at least a non-empty Key field, and
// the handler writes the new value to the database before returning a status
// message.
func UpdateSetting(c *gin.Context) {
	var s models.Setting
	if err := c.BindJSON(&s); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if s.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key is required"})
		return
	}
	if err := database.SetSettingValue(s.Key, s.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Saved"})
}
