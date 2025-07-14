package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	ModelID              uint    `json:"modelId"`
	VersionID            int     `gorm:"uniqueIndex" json:"versionId"`
	Name                 string  `json:"name"`
	BaseModel            string  `json:"baseModel"`
	EarlyAccessTimeFrame int     `json:"earlyAccessTimeFrame"`
	SizeKB               float64 `json:"sizeKB"`
	TrainedWords         string  `json:"trainedWords"`
	Nsfw                 bool    `json:"nsfw"`
	Type                 string  `json:"type"`
	Tags                 string  `json:"tags"`
	Description          string  `json:"description"`
	Mode                 string  `json:"mode"`
	ModelURL             string  `json:"modelUrl"`
	ImagePath            string  `json:"imagePath"`
	FilePath             string  `json:"filePath"`

	Images []VersionImage `json:"images"`
}
