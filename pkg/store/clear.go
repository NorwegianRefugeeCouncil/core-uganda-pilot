package store

import (
	"github.com/nrc-no/core/pkg/sql/convert"
	"gorm.io/gorm"
)

func Clear(db *gorm.DB) error {

	var databases []*Database
	if err := db.Find(&databases).Error; err != nil {
		return err
	}

	for _, database := range databases {
		if err := convert.DeleteDatabaseIfExists(db, database.ID); err != nil {
			return err
		}
	}

	if err := db.Where("field_id = field_id").Delete(&Option{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Field{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Form{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Folder{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Database{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Organization{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&IdentityProvider{}).Error; err != nil {
		return err
	}

	return nil
}
