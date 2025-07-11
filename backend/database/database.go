package database

import (
	"model-manager/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("backend/models.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	database.AutoMigrate(&models.Model{}, &models.Version{})
	DB = database
}
