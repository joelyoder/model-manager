package api

import (
	"net/http"
	"os"
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
	root := "./backend/downloads"
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".safetensors" && ext != ".pt" {
			return nil
		}
		abs, _ := filepath.Abs(path)
		var count int64
		database.DB.Model(&models.Version{}).Where("file_path = ?", abs).Or("file_path = ?", path).Count(&count)
		if count == 0 {
			database.DB.Model(&models.Model{}).Where("file_path = ?", abs).Or("file_path = ?", path).Count(&count)
		}
		if count == 0 {
			orphans = append(orphans, info.Name())
		}
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"files": orphans})
}
