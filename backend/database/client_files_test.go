package database

import (
	"model-manager/backend/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupClientFileTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.ClientFile{})
	DB = db
	return db
}

func TestResetAllPendingClientFiles(t *testing.T) {
	db := setupClientFileTestDB()

	// Seed data
	db.Create(&models.ClientFile{ClientID: "c1", Status: "pending"})
	db.Create(&models.ClientFile{ClientID: "c2", Status: "pending"})
	db.Create(&models.ClientFile{ClientID: "c1", Status: "installed"})

	err := ResetAllPendingClientFiles()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int64
	db.Model(&models.ClientFile{}).Where("status = ?", "pending").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 pending files, got %d", count)
	}

	db.Model(&models.ClientFile{}).Where("status = ?", "installed").Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 installed file, got %d", count)
	}
}

func TestResetPendingClientFilesForClient(t *testing.T) {
	db := setupClientFileTestDB()

	// Seed data
	db.Create(&models.ClientFile{ClientID: "c1", Status: "pending"})
	db.Create(&models.ClientFile{ClientID: "c1", Status: "pending"})
	db.Create(&models.ClientFile{ClientID: "c2", Status: "pending"})   // Should not be deleted
	db.Create(&models.ClientFile{ClientID: "c1", Status: "installed"}) // Should not be deleted

	err := ResetPendingClientFilesForClient("c1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var count int64
	db.Model(&models.ClientFile{}).Where("client_id = ? AND status = ?", "c1", "pending").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 pending files for c1, got %d", count)
	}

	db.Model(&models.ClientFile{}).Where("client_id = ? AND status = ?", "c2", "pending").Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 pending file for c2, got %d", count)
	}
}
