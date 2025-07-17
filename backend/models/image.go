package models

import "gorm.io/gorm"

// VersionImage represents a downloaded image for a model version
// including its local path and metadata extracted from CivitAI.
type VersionImage struct {
	gorm.Model
	VersionID uint   `gorm:"index" json:"versionId"`
	Path      string `json:"path"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Hash      string `json:"hash"`
	Meta      string `json:"meta"`
}
