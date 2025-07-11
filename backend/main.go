package main

import (
	"model-manager/backend/api"
	"model-manager/backend/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectDatabase()

	r := gin.Default()

	r.Static("/images", "./backend/images")
	r.Static("/downloads", "./backend/downloads")
	r.Static("/assets", "./frontend/dist/assets")

	// Model Versions
	r.GET("/api/model/:id/versions", api.GetModelVersions)
	r.POST("/api/sync/version/:versionId", api.SyncVersionByID)

	// API routes
	r.GET("/api/models", api.GetModels)
	r.POST("/api/sync", api.SyncCivitModels)
	r.POST("/api/sync/:id", api.SyncCivitModelByID)

	// Serve Vue SPA
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	r.Run(":8080")
}
