package api

import (
	"log"
	"net/http"
	"path/filepath"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// MigratePaths converts absolute paths in the database to relative paths based
// on the current configured roots. It is triggered manually by the user.
func MigratePaths(c *gin.Context) {
	modelRoot := database.GetModelPath()
	imageRoot := database.GetImagePath()

	log.Printf("Starting path migration. Model Root: %s, Image Root: %s", modelRoot, imageRoot)

	// Migrate Models
	var modelsList []models.Model
	if err := database.DB.Find(&modelsList).Error; err == nil {
		for _, m := range modelsList {
			updated := false
			if m.FilePath != "" && filepath.IsAbs(m.FilePath) {
				newPath := MakeRelativePath(m.FilePath, modelRoot)
				if newPath != m.FilePath {
					m.FilePath = newPath
					updated = true
				}
			}
			if m.ImagePath != "" && filepath.IsAbs(m.ImagePath) {
				newPath := MakeRelativePath(m.ImagePath, imageRoot)
				if newPath != m.ImagePath {
					m.ImagePath = newPath
					updated = true
				}
			}
			if updated {
				database.DB.Save(&m)
			}
		}
	}

	// Migrate Versions
	var versions []models.Version
	if err := database.DB.Find(&versions).Error; err == nil {
		for _, v := range versions {
			updated := false
			if v.FilePath != "" && filepath.IsAbs(v.FilePath) {
				newPath := MakeRelativePath(v.FilePath, modelRoot)
				if newPath != v.FilePath {
					v.FilePath = newPath
					updated = true
				}
			}
			if v.ImagePath != "" && filepath.IsAbs(v.ImagePath) {
				newPath := MakeRelativePath(v.ImagePath, imageRoot)
				if newPath != v.ImagePath {
					v.ImagePath = newPath
					updated = true
				}
			}
			if updated {
				database.DB.Save(&v)
			}
		}
	}

	// Migrate VersionImages
	var images []models.VersionImage
	if err := database.DB.Find(&images).Error; err == nil {
		for _, img := range images {
			if img.Path != "" && filepath.IsAbs(img.Path) {
				newPath := MakeRelativePath(img.Path, imageRoot)
				if newPath != img.Path {
					img.Path = newPath
					database.DB.Save(&img)
				}
			}
		}
	}

	log.Println("Path migration complete")
	c.JSON(http.StatusOK, gin.H{"message": "Path migration complete"})
}

// NormalizeSlashes is a helper to ensure we store forward slashes in DB
// even on Windows, to make the DB portable to Linux.
func NormalizeSlashes(path string) string {
	return filepath.ToSlash(path)
}
