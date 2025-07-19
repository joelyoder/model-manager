package api

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetOrphanFiles scans the downloads folder for .safetensors or .pt files
// that are not referenced by any model or version record. It returns a list
// of file names that have no associated record.
func GetOrphanFiles(c *gin.Context) {
	var orphans []string
	root := filepath.Join("backend", "downloads")
	rootAbs, _ := filepath.Abs(root)

	filepath.WalkDir(rootAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".safetensors" && ext != ".pt" {
			return nil
		}

		abs := filepath.Clean(path)
		norm := filepath.ToSlash(abs)

		var count int64
		database.DB.Model(&models.Version{}).Where("file_path = ?", norm).Or("file_path = ?", abs).Count(&count)
		if count == 0 {
			database.DB.Model(&models.Model{}).Where("file_path = ?", norm).Or("file_path = ?", abs).Count(&count)
		}
		if count == 0 {
			rel, relErr := filepath.Rel(rootAbs, abs)
			if relErr != nil {
				rel = d.Name()
			}
			orphans = append(orphans, rel)
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"files": orphans})
}
