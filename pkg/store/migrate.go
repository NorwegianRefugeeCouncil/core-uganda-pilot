package store

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Database{},
		&Form{},
		&Field{},
		&Folder{},
		&Organization{},
		&IdentityProvider{})
}
