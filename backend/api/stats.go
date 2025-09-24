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
// query parameters category, baseModel, modelType, and nsfw filter the
// dataset before aggregating counts. The handler performs read-only queries and
// responds with JSON metrics without mutating database state.
func GetStats(c *gin.Context) {
	category := strings.ToLower(strings.TrimSpace(c.Query("category")))
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")
	nsfwFilter := strings.ToLower(strings.TrimSpace(c.Query("nsfw")))
	if nsfwFilter != "" && nsfwFilter != "nsfw" && nsfwFilter != "non" {
		nsfwFilter = ""
	}

	versionFiltersActive := category != "" || baseModel != "" || modelType != "" || nsfwFilter != ""

	query := database.DB.Model(&models.Version{}).
		Joins("JOIN models ON models.id = versions.model_id").
		Where("models.deleted_at IS NULL").
		Where("versions.deleted_at IS NULL")

	if nsfwFilter == "non" {
		query = query.Where("versions.nsfw = ?", false)
	}
	if nsfwFilter == "nsfw" {
		query = query.Where("versions.nsfw = ?", true)
	}
	if baseModel != "" {
		query = query.Where("versions.base_model = ?", baseModel)
	}
	if modelType != "" {
		query = query.Where("versions.type = ?", modelType)
	}
	if category != "" {
		query = query.Where("LOWER(versions.tags) LIKE ?", "%"+category+"%")
	}

	var versions []models.Version
	if err := query.
		Select([]string{"versions.id", "versions.model_id", "versions.base_model", "versions.nsfw", "versions.type", "versions.tags"}).
		Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}

	baseCounts := make(map[string]int64)
	typeCounts := make(map[string]int64)
	categoryCounts := make(map[string]int64)

	filteredVersions := make([]models.Version, 0, len(versions))
	for _, v := range versions {
		if category != "" && !tagMatchesCategory(v.Tags, category) {
			continue
		}
		filteredVersions = append(filteredVersions, v)
	}

	var nsfwCount int64
	var safeCount int64
	modelSeen := make(map[uint]struct{})
	for _, v := range filteredVersions {
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
		modelSeen[v.ModelID] = struct{}{}
	}

	totalVersions := int64(len(filteredVersions))

	var totalModels int64
	if versionFiltersActive {
		totalModels = int64(len(modelSeen))
	} else {
		if err := database.DB.Model(&models.Model{}).Where("deleted_at IS NULL").Count(&totalModels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
			return
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
