package models

import "gorm.io/gorm"

type Collection struct {
	gorm.Model
	Name        string    `json:"name" gorm:"index"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions" gorm:"many2many:collection_versions;"`
}
