package models

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	CivitID     int `gorm:"uniqueIndex"`
	Name        string
	Type        string
	Nsfw        bool `gorm:"column:nsfw"`
	Description string
	CreatedAt   string
	UpdatedAt   string

	// Local paths
	ImagePath string
	FilePath  string
}
