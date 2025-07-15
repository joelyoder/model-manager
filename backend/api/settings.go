package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetSettings returns all settings as a key/value map.
func GetSettings(c *gin.Context) {
	var all []models.Setting
	database.DB.Find(&all)
	out := make(map[string]string)
	for _, s := range all {
		out[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, out)
}

// UpdateSettings updates or creates settings based on the provided map.
func UpdateSettings(c *gin.Context) {
	var input map[string]string
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	for k, v := range input {
		database.SetSetting(k, v)
	}
	c.JSON(http.StatusOK, gin.H{"message": "settings saved"})
}
