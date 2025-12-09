package database

import "model-manager/backend/models"

// GetSettingValue returns the value for the given key or empty string if not found.
// GetSettingValue returns the value for the given key or empty string if not found.
func GetSettingValue(key string) string {
	var s models.Setting
	// Use struct-based query to avoid SQL keyword issues
	result := DB.Where(&models.Setting{Key: key}).First(&s)
	if result.Error == nil {
		// Log success for debugging (verbose, but necessary right now)
		// fmt.Println("GetSettingValue", key, "found:", s.Value)
		return s.Value
	}
	// Log missing key only if it's expected to be there (ignore some like migration keys if common)
	if key == "model_path" || key == "image_path" {
		// fmt.Println("GetSettingValue", key, "NOT FOUND, using default")
	}
	return ""
}

// SetSettingValue creates or updates the setting with the specified key.
func SetSettingValue(key, value string) error {
	var s models.Setting
	if err := DB.Where(&models.Setting{Key: key}).First(&s).Error; err == nil {
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
