package api

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// isAbsolutePath detects absolute paths cross-platform.
// On Linux, filepath.IsAbs won't recognize Windows paths like "C:\..." and vice versa.
// This helper detects both formats regardless of the current OS.
func isAbsolutePath(path string) bool {
	// Standard check for current OS
	if filepath.IsAbs(path) {
		return true
	}
	// Check for Windows-style absolute paths (e.g., "C:\..." or "D:/...")
	if len(path) >= 2 && path[1] == ':' {
		return true
	}
	// Check for Unix-style absolute paths (starts with /)
	if strings.HasPrefix(path, "/") {
		return true
	}
	return false
}

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
			if m.FilePath != "" && isAbsolutePath(m.FilePath) {
				newPath := MakeRelativePath(m.FilePath, modelRoot)
				if newPath != m.FilePath {
					m.FilePath = newPath
					updated = true
				}
			}
			if m.ImagePath != "" && isAbsolutePath(m.ImagePath) {
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
			if v.FilePath != "" && isAbsolutePath(v.FilePath) {
				newPath := MakeRelativePath(v.FilePath, modelRoot)
				if newPath != v.FilePath {
					v.FilePath = newPath
					updated = true
				}
			}
			if v.ImagePath != "" && isAbsolutePath(v.ImagePath) {
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
			if img.Path != "" && isAbsolutePath(img.Path) {
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
