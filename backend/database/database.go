package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"model-manager/backend/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("backend/models.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	database.AutoMigrate(&models.Model{}, &models.Version{}, &models.VersionImage{}, &models.Setting{})
	DB = database
}
