package api

import (
	"net/http"
	"strings"
	"time"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// ExportModels returns all model data in a JSON format compatible with ImportModels.
func ExportModels(c *gin.Context) {
	var versions []models.Version
	database.DB.Find(&versions)

	records := make([]ImportRecord, 0, len(versions))
	for _, v := range versions {
		var model models.Model
		database.DB.First(&model, v.ModelID)

		// Convert tags string into slice
		groups := []string{}
		if model.Tags != "" {
			for _, t := range strings.Split(model.Tags, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					groups = append(groups, t)
				}
			}
		}
		if model.Nsfw {
			groups = append(groups, "nsfw")
		}

		created := 0.0
		if v.CivitCreatedAt != "" {
			if t, err := time.Parse(time.RFC3339, v.CivitCreatedAt); err == nil {
				created = float64(t.Unix())
			}
		}

		name := model.Name
		if v.Name != "" {
			name = model.Name + " (" + v.Name + ")"
		}

		records = append(records, ImportRecord{
			Name:            name,
			BaseModel:       v.BaseModel,
			ModelType:       v.Type,
			DownloadURL:     v.DownloadURL,
			URL:             v.ModelURL,
			PreviewURL:      "",
			Description:     v.Description,
			PositivePrompts: v.TrainedWords,
			SHA256Hash:      v.SHA256,
			CreatedAt:       created,
			Groups:          groups,
			Location:        v.FilePath,
		})
	}

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=\"model_export.json\"")
	c.JSON(http.StatusOK, records)
}
