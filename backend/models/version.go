package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	ModelID              uint    `gorm:"index" json:"modelId"`
	ParentModel          Model   `json:"model" gorm:"foreignKey:ModelID"`
	VersionID            int     `gorm:"uniqueIndex" json:"versionId"`
	Name                 string  `gorm:"index" json:"name"`
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
	CivitCreatedAt       string  `json:"createdAt"`
	CivitUpdatedAt       string  `json:"updatedAt"`
	SHA256               string  `json:"sha256"`
	DownloadURL          string  `json:"downloadUrl"`
	ImagePath            string  `json:"imagePath"`
	FilePath             string  `json:"filePath"`

	Images []VersionImage `json:"images"`

	// ClientStatus is a calculated field for the default client, not stored in DB
	ClientStatus string `gorm:"-" json:"clientStatus"`
}
