package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetSettings returns all settings.
func GetSettings(c *gin.Context) {
	settings, err := database.GetAllSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

// UpdateSetting updates or creates a setting.
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
