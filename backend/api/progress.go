package api

import "github.com/gin-gonic/gin"

func GetDownloadProgress(c *gin.Context) {
	c.JSON(200, gin.H{"progress": CurrentDownloadProgress})
}
