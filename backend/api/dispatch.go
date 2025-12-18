package api

import (
	"log"
	"net/http"
	"path/filepath"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

type DispatchRequest struct {
	Action         string `json:"action"` // "download", "delete"
	URL            string `json:"url"`
	Filename       string `json:"filename"`
	Subdirectory   string `json:"subdirectory"`
	ModelVersionID uint   `json:"model_version_id"`
	ClientID       string `json:"client_id"`
}

func DispatchRemote(c *gin.Context) {
	var req DispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ClientID == "" {
		// Fallback for single-client setups? Or error?
		// For now, error.
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id is required"})
		return
	}

	// Lookup Version to get local FilePath
	var version models.Version
	if err := database.DB.Preload("ParentModel").First(&version, req.ModelVersionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model version not found"})
		return
	}

	if version.FilePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Model version has no local file"})
		return
	}

	// Update DB state
	if req.Action == "download" {
		var cf models.ClientFile
		// Check if exists using Find to avoid RecordNotFound log
		result := database.DB.Where("client_id = ? AND model_version_id = ?", req.ClientID, req.ModelVersionID).Limit(1).Find(&cf)
		if result.RowsAffected == 0 {
			cf = models.ClientFile{
				ClientID:       req.ClientID,
				ModelVersionID: req.ModelVersionID,
			}
		}
		cf.Status = "pending"
		database.DB.Save(&cf)

		// Construct clean relative path for URL
		// We use ParentModel.Type as subdirectory if available, or Version.Type
		subdir := version.ParentModel.Type
		if subdir == "" {
			subdir = version.Type
		}
		if subdir == "" {
			subdir = "Other" // Fallback
		}

		filename := filepath.Base(version.FilePath)
		relativePath := filepath.ToSlash(filepath.Join(subdir, filename))
		req.URL = "/downloads/" + relativePath

		log.Printf("Dispatching download for model %d to client %s (local_url=%s)", req.ModelVersionID, req.ClientID, req.URL)
	} else if req.Action == "delete" {
		// Maybe set to "removing"? Or just let client confirm deletion?
		// Prompt says: If Action is "delete": Update ClientFile record to remove the entry (or set to "removed").
		// I'll wait for confirmation to remove, but maybe set status to 'pending_delete' if I wanted to be fancy.
		// For now, I'll delete it immediately from DB or wait? User says "Update... to remove the entry".
		// I'll delete it now, or maybe better to wait for callback?
		// "Action is delete: Update ClientFile record to remove the entry"
		database.DB.Delete(&models.ClientFile{}, "client_id = ? AND model_version_id = ?", req.ClientID, req.ModelVersionID)

		// Use filepath from DB to ensure correct file is deleted
		subdir := version.ParentModel.Type
		if subdir == "" {
			subdir = version.Type
		}
		if subdir == "" {
			subdir = "Other"
		}
		filename := filepath.Base(version.FilePath)
		// Client expects filename to be the relative path
		req.Filename = filepath.ToSlash(filepath.Join(subdir, filename))
		req.Subdirectory = ""

		log.Printf("Dispatching delete for model %d to client %s", req.ModelVersionID, req.ClientID)
	}

	// Send to WebSocket
	err := SendToClient(req.ClientID, req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Client not connected", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "dispatched"})
}
