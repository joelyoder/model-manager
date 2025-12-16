package models

import "gorm.io/gorm"

type ClientFile struct {
	gorm.Model
	ClientID       string `gorm:"index" json:"clientId"`
	ModelVersionID uint   `gorm:"index" json:"modelVersionId"`
	Status         string `json:"status"` // "pending", "installed"
}
