package models

import "gorm.io/gorm"

type Version struct {
	gorm.Model
	ModelID              uint
	VersionID            int    `gorm:"uniqueIndex"`
	Name                 string
	BaseModel            string
	CreatedAt            string
	EarlyAccessTimeFrame int
	SizeKB               float64
	TrainedWords         string
	ImagePath            string
	FilePath             string
}
