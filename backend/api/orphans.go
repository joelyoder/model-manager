package api

import (
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

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
			log.Printf("failed to get abs path for %s: %v", p, err)
			continue
		}
		key := abs
		if runtime.GOOS == "windows" {
			key = strings.ToLower(key)
		}
		dbFiles[key] = struct{}{}
		log.Printf("db file: %s", key)
	}

	var orphans []string
	root := filepath.Join("backend", "downloads")
	log.Printf("walking downloads directory: %s", root)
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("walk error on %s: %v", path, err)
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".safetensors" && ext != ".pt" {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			log.Printf("failed to get abs path for %s: %v", path, err)
			return nil
		}
		key := abs
		if runtime.GOOS == "windows" {
			key = strings.ToLower(key)
		}
		if _, exists := dbFiles[key]; !exists {
			log.Printf("orphaned file: %s", abs)
			orphans = append(orphans, abs)
		} else {
			log.Printf("file referenced in db: %s", abs)
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"orphans": orphans})
}
