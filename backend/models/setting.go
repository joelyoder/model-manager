package models

import "gorm.io/gorm"

// Setting represents a key-value pair application configuration.
type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex" json:"key"`
	Value string `json:"value"`
}
