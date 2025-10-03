package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDownloadProgress reports the global download progress percentage for
// long-running model downloads. The endpoint accepts no parameters and simply
// reads the package-level CurrentDownloadProgress value to return JSON.
func GetDownloadProgress(c *gin.Context) {
	c.JSON(200, gin.H{"progress": CurrentDownloadProgress})
}

// CancelDownload stops the active model download, if any, and cleans up any
// partially written files. It responds with a success message regardless of
// whether a download was in progress so the action remains idempotent.
func CancelDownload(c *gin.Context) {
	cancelled, err := CancelActiveDownload()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel download"})
		return
	}
	if cancelled {
		c.JSON(http.StatusOK, gin.H{"message": "Download cancelled"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "No active download"})
}
