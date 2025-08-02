package api

import (
	"net/http"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

type countResult struct {
	Key   string
	Count int64
}

// GetStats returns summary statistics about the models in the database.
func GetStats(c *gin.Context) {
	var total int64
	database.DB.Model(&models.Model{}).Count(&total)

	var typeResults []countResult
	database.DB.Model(&models.Model{}).
		Select("type as key, count(*) as count").
		Group("type").Scan(&typeResults)

	var baseResults []countResult
	database.DB.Model(&models.Version{}).
		Select("base_model as key, count(*) as count").
		Group("base_model").Scan(&baseResults)

	var nsfwRows []struct {
		Nsfw  bool
		Count int64
	}
	database.DB.Model(&models.Model{}).
		Select("nsfw, count(*) as count").
		Group("nsfw").Scan(&nsfwRows)

	nsfwCount := int64(0)
	safeCount := int64(0)
	for _, r := range nsfwRows {
		if r.Nsfw {
			nsfwCount += r.Count
		} else {
			safeCount += r.Count
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"totalModels":     total,
		"typeCounts":      typeResults,
		"baseModelCounts": baseResults,
		"nsfwCount":       nsfwCount,
		"nonNsfwCount":    safeCount,
	})
}
