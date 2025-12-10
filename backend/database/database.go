package database

import (
	"log"
	"os"

	"model-manager/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	path := os.Getenv("MODELS_DB_PATH")
	if path == "" {
		path = "backend/models.db"
	}
	log.Printf("Connecting to database at: %s", path)
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	database.AutoMigrate(&models.Model{}, &models.Version{}, &models.VersionImage{}, &models.Setting{})
	DB = database

	if err := applyMigrations(database); err != nil {
		log.Printf("failed to run database migrations: %v", err)
	}

	// Log configured paths at startup for debugging
	log.Printf("Configured model_path: '%s'", GetModelPath())
	log.Printf("Configured image_path: '%s'", GetImagePath())
}
