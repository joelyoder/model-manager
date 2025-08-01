package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetStats returns library statistics including counts and breakdowns by
// type, base model and NSFW status.
func GetStats(c *gin.Context) {
	var total int64
	database.DB.Model(&models.Model{}).Count(&total)

	type TypeCount struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	var typeCounts []TypeCount
	database.DB.Model(&models.Model{}).
		Select("type as type, COUNT(*) as count").
		Group("type").
		Order("type").
		Find(&typeCounts)

	type BaseCount struct {
		BaseModel string `json:"baseModel"`
		Count     int64  `json:"count"`
	}
	var baseCounts []BaseCount
	database.DB.Model(&models.Version{}).
		Select("base_model as base_model, COUNT(*) as count").
		Group("base_model").
		Order("base_model").
		Find(&baseCounts)

	var nsfwCount int64
	database.DB.Model(&models.Model{}).Where("nsfw = 1").Count(&nsfwCount)
	sfwCount := total - nsfwCount

	c.JSON(http.StatusOK, gin.H{
		"totalModels": total,
		"types":       typeCounts,
		"baseModels":  baseCounts,
		"nsfwCount":   nsfwCount,
		"sfwCount":    sfwCount,
	})
}
