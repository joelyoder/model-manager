package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"model-manager/backend/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	path := os.Getenv("MODELS_DB_PATH")
	if path == "" {
		path = "backend/models.db"
	}
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	database.AutoMigrate(&models.Model{}, &models.Version{}, &models.VersionImage{}, &models.Setting{})
	DB = database
}
