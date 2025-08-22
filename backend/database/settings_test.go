package database

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"model-manager/backend/models"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { sqlDB.Close() })

	if err := db.AutoMigrate(&models.Setting{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	DB = db
	return db
}

func TestGetSettingValue(t *testing.T) {
	setupTestDB(t)

	// Seed a setting
	if err := DB.Create(&models.Setting{Key: "theme", Value: "dark"}).Error; err != nil {
		t.Fatalf("failed to create setting: %v", err)
	}

	if value := GetSettingValue("theme"); value != "dark" {
		t.Fatalf("expected 'dark', got %q", value)
	}

	if value := GetSettingValue("missing"); value != "" {
		t.Fatalf("expected empty string for missing key, got %q", value)
	}
}

func TestSetSettingValue(t *testing.T) {
	setupTestDB(t)

	// Create new key
	if err := SetSettingValue("language", "en"); err != nil {
		t.Fatalf("failed to set new setting: %v", err)
	}
	var s models.Setting
	if err := DB.First(&s, "key = ?", "language").Error; err != nil {
		t.Fatalf("setting not created: %v", err)
	}
	if s.Value != "en" {
		t.Fatalf("expected 'en', got %q", s.Value)
	}

	// Update existing key
	if err := SetSettingValue("language", "fr"); err != nil {
		t.Fatalf("failed to update setting: %v", err)
	}
	var count int64
	DB.Model(&models.Setting{}).Where("key = ?", "language").Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 record, got %d", count)
	}
	if err := DB.First(&s, "key = ?", "language").Error; err != nil {
		t.Fatalf("setting not found after update: %v", err)
	}
	if s.Value != "fr" {
		t.Fatalf("expected 'fr', got %q", s.Value)
	}
}

func TestGetAllSettings(t *testing.T) {
	setupTestDB(t)

	if err := DB.Create(&models.Setting{Key: "a", Value: "1"}).Error; err != nil {
		t.Fatalf("failed to create setting a: %v", err)
	}
	if err := DB.Create(&models.Setting{Key: "b", Value: "2"}).Error; err != nil {
		t.Fatalf("failed to create setting b: %v", err)
	}

	settings, err := GetAllSettings()
	if err != nil {
		t.Fatalf("GetAllSettings returned error: %v", err)
	}
	if len(settings) != 2 {
		t.Fatalf("expected 2 settings, got %d", len(settings))
	}

	expected := map[string]string{"a": "1", "b": "2"}
	for _, s := range settings {
		if expected[s.Key] != s.Value {
			t.Errorf("unexpected value for key %s: %s", s.Key, s.Value)
		}
		delete(expected, s.Key)
	}
	if len(expected) != 0 {
		t.Fatalf("missing settings: %v", expected)
	}
}
