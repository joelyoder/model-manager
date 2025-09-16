package api

import "github.com/gin-gonic/gin"

// GetDownloadProgress reports the global download progress percentage for
// long-running model downloads. The endpoint accepts no parameters and simply
// reads the package-level CurrentDownloadProgress value to return JSON.
func GetDownloadProgress(c *gin.Context) {
	c.JSON(200, gin.H{"progress": CurrentDownloadProgress})
}
