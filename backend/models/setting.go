package models

import "gorm.io/gorm"

// Setting represents a key/value configuration stored in the database
// for application level configuration like API keys and folder paths.
type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex"`
	Value string
}

const (
	// SettingAPIKey stores the Civitai API key
	SettingAPIKey = "api_key"
	// SettingImagesPath stores the root folder for downloaded images
	SettingImagesPath = "images_path"
	// SettingModelsPath stores the root folder for downloaded model files
	SettingModelsPath = "models_path"
)
