package database

import "model-manager/backend/models"

// GetSettingValue returns the value for the given key or empty string if not found.
func GetSettingValue(key string) string {
	var s models.Setting
	if err := DB.First(&s, "key = ?", key).Error; err == nil {
		return s.Value
	}
	return ""
}

// SetSettingValue creates or updates the setting with the specified key.
func SetSettingValue(key, value string) error {
	var s models.Setting
	if err := DB.First(&s, "key = ?", key).Error; err == nil {
		s.Value = value
		return DB.Save(&s).Error
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
