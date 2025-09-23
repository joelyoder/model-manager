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

// GetStats aggregates counts about the stored models and their versions and
// returns totals grouped by model type, base model, and NSFW flag. Optional
// query parameters category, baseModel, modelType, and hideNsfw filter the
// dataset before aggregating counts. The handler performs read-only queries and
// responds with JSON metrics without mutating database state.
func GetStats(c *gin.Context) {
	category := strings.ToLower(strings.TrimSpace(c.Query("category")))
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")
	hideNsfw := c.Query("hideNsfw") == "1"

	var modelsList []models.Model
	if err := database.DB.Preload("Versions").Find(&modelsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}

	versionFiltersActive := category != "" || baseModel != "" || modelType != ""

	baseCounts := make(map[string]int64)
	typeCounts := make(map[string]int64)
	categoryCounts := make(map[string]int64)

	type modelAggregate struct {
		Model    models.Model
		Versions []models.Version
	}

	aggregates := make([]modelAggregate, 0, len(modelsList))
	includedVersions := make([]models.Version, 0)

	for _, m := range modelsList {
		if hideNsfw && m.Nsfw {
			continue
		}

		matchingVersions := make([]models.Version, 0, len(m.Versions))
		for _, v := range m.Versions {
			if hideNsfw && v.Nsfw {
				continue
			}
			if baseModel != "" && v.BaseModel != baseModel {
				continue
			}
			if modelType != "" && v.Type != modelType {
				continue
			}
			if category != "" && !tagMatchesCategory(v.Tags, category) {
				continue
			}

			matchingVersions = append(matchingVersions, v)
		}

		includeModel := len(matchingVersions) > 0 || (len(m.Versions) == 0 && !versionFiltersActive)
		if !includeModel {
			continue
		}

		aggregates = append(aggregates, modelAggregate{Model: m, Versions: matchingVersions})
		includedVersions = append(includedVersions, matchingVersions...)
	}

	var nsfwCount int64
	var safeCount int64
	for _, v := range includedVersions {
		baseCounts[v.BaseModel]++
		typeCounts[v.Type]++

		cats := extractCategories(v.Tags)
		if len(cats) == 0 {
			categoryCounts[uncategorizedLabel]++
		} else {
			for _, cat := range cats {
				categoryCounts[cat]++
			}
		}

		if v.Nsfw {
			nsfwCount++
		} else {
			safeCount++
		}
	}

	totalModels := int64(len(aggregates))
	totalVersions := int64(len(includedVersions))

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

	categoryResults := make([]countResult, 0, len(categoryCounts))
	for k, v := range categoryCounts {
		categoryResults = append(categoryResults, countResult{Key: k, Count: v})
	}
	sort.Slice(categoryResults, func(i, j int) bool {
		return categoryResults[i].Key < categoryResults[j].Key
	})

	c.JSON(http.StatusOK, gin.H{
		"totalModels":     totalModels,
		"totalVersions":   totalVersions,
		"typeCounts":      typeResults,
		"baseModelCounts": baseResults,
		"categoryCounts":  categoryResults,
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

const uncategorizedLabel = "Uncategorized"

var categoryLookup = map[string]string{
	"character":  "Character",
	"style":      "Style",
	"concept":    "Concept",
	"clothing":   "Clothing",
	"base model": "Base Model",
	"poses":      "Poses",
	"background": "Background",
	"tool":       "Tool",
	"vehicle":    "Vehicle",
	"buildings":  "Buildings",
	"objects":    "Objects",
	"assets":     "Assets",
	"animal":     "Animal",
	"action":     "Action",
}

func extractCategories(tags string) []string {
	seen := make(map[string]struct{})
	for _, t := range strings.Split(tags, ",") {
		normalized := strings.TrimSpace(strings.ToLower(t))
		if normalized == "" {
			continue
		}
		if display, ok := categoryLookup[normalized]; ok {
			seen[display] = struct{}{}
		}
	}
	if len(seen) == 0 {
		return nil
	}
	results := make([]string, 0, len(seen))
	for cat := range seen {
		results = append(results, cat)
	}
	return results
}
