package database

import (
	"model-manager/backend/models"
)

// GetSetting returns the value for the specified key or empty string if not set.
func GetSetting(key string) string {
	var s models.Setting
	if err := DB.Where("key = ?", key).First(&s).Error; err == nil {
		return s.Value
	}
	return ""
}

// SetSetting creates or updates a key/value pair in the settings table.
func SetSetting(key, value string) {
	var s models.Setting
	if err := DB.Where("key = ?", key).First(&s).Error; err == nil {
		s.Value = value
		DB.Save(&s)
	} else {
		DB.Create(&models.Setting{Key: key, Value: value})
	}
}

// GetImagesPath returns the configured images folder or the default path.
func GetImagesPath() string {
	p := GetSetting(models.SettingImagesPath)
	if p == "" {
		p = "./backend/images"
	}
	return p
}

// GetModelsPath returns the configured models folder or the default path.
func GetModelsPath() string {
	p := GetSetting(models.SettingModelsPath)
	if p == "" {
		p = "./backend/downloads"
	}
	return p
}

// GetAPIKey returns the configured Civitai API key.
func GetAPIKey() string {
	return GetSetting(models.SettingAPIKey)
}
