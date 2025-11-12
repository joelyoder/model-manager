package database

import (
	"model-manager/backend/models"

	"gorm.io/gorm"
)

const modelWeightMigrationKey = "migration:model-weight-defaulted"

func applyMigrations(db *gorm.DB) error {
	if err := backfillModelWeights(db); err != nil {
		return err
	}
	return nil
}

func backfillModelWeights(db *gorm.DB) error {
	if GetSettingValue(modelWeightMigrationKey) == "1" {
		return nil
	}

	if err := db.Model(&models.Model{}).
		Where("weight <= 0 OR weight IS NULL").
		Update("weight", 1).Error; err != nil {
		return err
	}

	if err := SetSettingValue(modelWeightMigrationKey, "1"); err != nil {
		return err
	}
	return nil
}
