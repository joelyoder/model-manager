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
		apiGroup.GET("/models/:id", api.GetModel)
		apiGroup.DELETE("/models/:id", api.DeleteModel)
		apiGroup.POST("/sync", api.SyncCivitModels)
		apiGroup.POST("/sync/:id", api.SyncCivitModelByID)
		apiGroup.POST("/sync/version/:versionId", api.SyncVersionByID)
		apiGroup.GET("/model/:id/versions", api.GetModelVersions)
		apiGroup.DELETE("/versions/:id", api.DeleteVersion)
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
