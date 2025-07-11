package models

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	CivitID     int    `gorm:"uniqueIndex" json:"civitId"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Tags        string `json:"tags"`
	Nsfw        bool   `gorm:"column:nsfw" json:"nsfw"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`

	// Local paths
	ImagePath string `json:"imagePath"`
	FilePath  string `json:"filePath"`

	ImageWidth  int `json:"imageWidth"`
	ImageHeight int `json:"imageHeight"`

	Versions []Version `json:"versions"`
}
