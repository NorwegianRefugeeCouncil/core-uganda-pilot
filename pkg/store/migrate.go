package store

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Database{},
		&Form{},
		&Field{},
		&Option{},
		&Folder{},
		&Organization{},
		&IdentityProvider{})
}
