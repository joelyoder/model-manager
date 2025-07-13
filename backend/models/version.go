package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	ModelID              uint    `json:"modelId"`
	VersionID            int     `gorm:"uniqueIndex" json:"versionId"`
	Name                 string  `json:"name"`
	BaseModel            string  `json:"baseModel"`
	CreatedAt            string  `json:"createdAt"`
	EarlyAccessTimeFrame int     `json:"earlyAccessTimeFrame"`
	SizeKB               float64 `json:"sizeKB"`
	TrainedWords         string  `json:"trainedWords"`
	ModelURL             string  `json:"modelUrl"`
	ImagePath            string  `json:"imagePath"`
	FilePath             string  `json:"filePath"`
}
