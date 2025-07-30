package main

import (
	"log"
	"os"

	"model-manager/backend/api"
	"model-manager/backend/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectDatabase()

	r := gin.Default()
	r.SetTrustedProxies(nil) // safe for local dev

	// Serve static assets
	r.Static("/images", "./backend/images")
	r.Static("/downloads", "./backend/downloads")
	r.Static("/assets", "./frontend/dist/assets")

	// API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/models", api.GetModels)
		apiGroup.GET("/models/count", api.GetModelsCount)
		apiGroup.GET("/base-models", api.GetBaseModels)
		apiGroup.GET("/models/:id", api.GetModel)
		apiGroup.PUT("/models/:id", api.UpdateModel)
		apiGroup.DELETE("/models/:id", api.DeleteModel)
		apiGroup.POST("/sync", api.SyncCivitModels)
		apiGroup.POST("/sync/:id", api.SyncCivitModelByID)
		apiGroup.POST("/sync/version/:versionId", api.SyncVersionByID)
		apiGroup.GET("/download/progress", api.GetDownloadProgress)
		apiGroup.GET("/model/:id/versions", api.GetModelVersions)
		apiGroup.GET("/versions/:id", api.GetVersion)
		apiGroup.PUT("/versions/:id", api.UpdateVersion)
		apiGroup.POST("/versions/:id/refresh", api.RefreshVersion)
		apiGroup.POST("/versions/:id/main-image/:imageId", api.SetVersionMainImage)
		apiGroup.DELETE("/versions/:id", api.DeleteVersion)
		apiGroup.POST("/import", api.ImportModels)
		apiGroup.POST("/import-db", api.ImportDatabase)
		apiGroup.GET("/export", api.ExportModels)
		apiGroup.GET("/settings", api.GetSettings)
		apiGroup.POST("/settings", api.UpdateSetting)
	}

	// Vue SPA fallback for all other routes (no wildcard)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	r.Run(":" + port)
}
