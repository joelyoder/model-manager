package api

import (
	"net/http"

	"model-manager/backend/database"

	"github.com/gin-gonic/gin"
)

// GetDuplicateFilePaths searches for versions sharing the same file path
func GetDuplicateFilePaths(c *gin.Context) {
	var rows []struct {
		ModelID     uint
		ModelName   string
		VersionID   uint
		VersionName string
		FilePath    string
	}
	database.DB.Table("versions").
		Select("versions.model_id, models.name as model_name, versions.id as version_id, versions.name as version_name, versions.file_path").
		Joins("JOIN models ON models.id = versions.model_id").
		Where("versions.file_path <> ''").
		Scan(&rows)

	type dupVersion struct {
		ModelID     uint   `json:"modelId"`
		ModelName   string `json:"modelName"`
		VersionID   uint   `json:"versionId"`
		VersionName string `json:"versionName"`
	}
	type dupFile struct {
		Path     string       `json:"path"`
		Versions []dupVersion `json:"versions"`
	}
	mp := make(map[string][]dupVersion)
	for _, r := range rows {
		v := dupVersion{
			ModelID:     r.ModelID,
			ModelName:   r.ModelName,
			VersionID:   r.VersionID,
			VersionName: r.VersionName,
		}
		mp[r.FilePath] = append(mp[r.FilePath], v)
	}
	var dups []dupFile
	for p, vs := range mp {
		if len(vs) > 1 {
			dups = append(dups, dupFile{Path: p, Versions: vs})
		}
	}
	c.JSON(http.StatusOK, gin.H{"duplicates": dups})
}
