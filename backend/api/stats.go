package api

import (
	"net/http"
	"sort"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

type countResult struct {
	Key   string
	Count int64
}

// GetStats aggregates counts about the stored models and returns totals,
// grouped counts by model type, base model, and NSFW flag. Optional query
// parameters category, baseModel, and modelType filter the dataset before
// aggregating counts. The handler performs read-only queries and responds with
// JSON metrics without mutating database state.
func GetStats(c *gin.Context) {
	category := strings.ToLower(strings.TrimSpace(c.Query("category")))
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")

	var modelsList []models.Model
	if err := database.DB.Preload("Versions").Find(&modelsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}

	filtersActive := category != "" || baseModel != "" || modelType != ""

	baseCounts := make(map[string]int64)
	typeCounts := make(map[string]int64)
	var total int64
	var nsfwCount int64
	var safeCount int64

	for _, m := range modelsList {
		includeModel := len(m.Versions) == 0 && !filtersActive

		for _, v := range m.Versions {
			if baseModel != "" && v.BaseModel != baseModel {
				continue
			}
			if modelType != "" && v.Type != modelType {
				continue
			}
			if category != "" && !tagMatchesCategory(v.Tags, category) {
				continue
			}

			includeModel = true

			key := v.BaseModel
			baseCounts[key]++
		}

		if !includeModel {
			continue
		}

		total++
		typeCounts[m.Type]++
		if m.Nsfw {
			nsfwCount++
		} else {
			safeCount++
		}
	}

	typeResults := make([]countResult, 0, len(typeCounts))
	for k, v := range typeCounts {
		typeResults = append(typeResults, countResult{Key: k, Count: v})
	}
	sort.Slice(typeResults, func(i, j int) bool {
		return typeResults[i].Key < typeResults[j].Key
	})

	baseResults := make([]countResult, 0, len(baseCounts))
	for k, v := range baseCounts {
		baseResults = append(baseResults, countResult{Key: k, Count: v})
	}
	sort.Slice(baseResults, func(i, j int) bool {
		return baseResults[i].Key < baseResults[j].Key
	})

	c.JSON(http.StatusOK, gin.H{
		"totalModels":     total,
		"typeCounts":      typeResults,
		"baseModelCounts": baseResults,
		"nsfwCount":       nsfwCount,
		"nonNsfwCount":    safeCount,
	})
}

func tagMatchesCategory(tags string, category string) bool {
	if category == "" {
		return true
	}
	for _, t := range strings.Split(tags, ",") {
		if strings.TrimSpace(strings.ToLower(t)) == category {
			return true
		}
	}
	return false
}
