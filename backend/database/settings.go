package database

import (
	"model-manager/backend/models"

	"gorm.io/gorm"
)

// GetSettingValue returns the value for the given key or empty string if not found.
func GetSettingValue(key string) string {
	var s models.Setting
	// Use struct-based query to avoid SQL keyword issues
	result := DB.Where(&models.Setting{Key: key}).First(&s)
	if result.Error == nil {
		return s.Value
	}
	return ""
}

// SetSettingValue creates or updates the setting with the specified key.
func SetSettingValue(key, value string) error {
	var s models.Setting
	// Use Unscoped to find even soft-deleted records to avoid unique index violation on creation
	if err := DB.Unscoped().Where(&models.Setting{Key: key}).First(&s).Error; err == nil {
		s.Value = value
		s.DeletedAt = gorm.DeletedAt{} // Restore if deleted
		return DB.Unscoped().Save(&s).Error
	}
	s.Key = key
	s.Value = value
	return DB.Create(&s).Error
}

// GetAllSettings returns all settings in the database.
func GetAllSettings() ([]models.Setting, error) {
	var settings []models.Setting
	if err := DB.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// GetModelPath returns the configured path for storing models.
// Defaults to "./backend/downloads" if not set.
func GetModelPath() string {
	if val := GetSettingValue("model_path"); val != "" {
		return val
	}
	return "./backend/downloads"
}

// GetImagePath returns the configured path for storing images.
// Defaults to "./backend/images" if not set.
func GetImagePath() string {
	if val := GetSettingValue("image_path"); val != "" {
		return val
	}
	return "./backend/images"
}
