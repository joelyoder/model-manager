package api

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetOrphanedFiles returns file paths in backend/downloads that are not referenced in the database.
func GetOrphanedFiles(c *gin.Context) {
	// Collect file paths from models and versions
	var modelPaths []string
	database.DB.Model(&models.Model{}).Where("file_path <> ''").Pluck("file_path", &modelPaths)
	var versionPaths []string
	database.DB.Model(&models.Version{}).Where("file_path <> ''").Pluck("file_path", &versionPaths)

	dbFiles := make(map[string]struct{})
	for _, p := range append(modelPaths, versionPaths...) {
		abs, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		dbFiles[abs] = struct{}{}
	}

	var orphans []string
	root := "./backend/downloads"
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			return nil
		}
		if _, exists := dbFiles[abs]; !exists {
			orphans = append(orphans, abs)
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"orphans": orphans})
}
