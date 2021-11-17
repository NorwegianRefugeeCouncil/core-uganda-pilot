package store

import (
	"gorm.io/gorm"
)

func Clear(db *gorm.DB) error {

	if err := db.Where("id = id").Delete(&CredentialIdentifier{}).Error; err != nil {
		return err
	}
	if err := db.Where("id = id").Delete(&Credential{}).Error; err != nil {
		return err
	}

	if err := db.Where("id = id").Delete(&Identity{}).Error; err != nil {
		return err
	}

	return nil
}
