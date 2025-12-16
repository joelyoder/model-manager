package api

import (
	"log"
	"net/http"

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

	// Update DB state
	if req.Action == "download" {
		var cf models.ClientFile
		// Check if exists
		res := database.DB.Where("client_id = ? AND model_version_id = ?", req.ClientID, req.ModelVersionID).First(&cf)
		if res.Error != nil {
			cf = models.ClientFile{
				ClientID:       req.ClientID,
				ModelVersionID: req.ModelVersionID,
			}
		}
		cf.Status = "pending"
		database.DB.Save(&cf)
		log.Printf("Dispatching download for model %d to client %s (status=pending)", req.ModelVersionID, req.ClientID)
	} else if req.Action == "delete" {
		// Maybe set to "removing"? Or just let client confirm deletion?
		// Prompt says: If Action is "delete": Update ClientFile record to remove the entry (or set to "removed").
		// I'll wait for confirmation to remove, but maybe set status to 'pending_delete' if I wanted to be fancy.
		// For now, I'll delete it immediately from DB or wait? User says "Update... to remove the entry".
		// I'll delete it now, or maybe better to wait for callback?
		// "Action is delete: Update ClientFile record to remove the entry"
		database.DB.Delete(&models.ClientFile{}, "client_id = ? AND model_version_id = ?", req.ClientID, req.ModelVersionID)
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
