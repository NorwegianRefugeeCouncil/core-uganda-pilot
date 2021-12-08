package mocks

import "gorm.io/gorm"

type DbFactory struct {
	db *gorm.DB
}

func (d *DbFactory) Get() (*gorm.DB, error) {
	return d.db, nil
}

// NewMockFactory returns a mock Factory
func NewMockFactory(db *gorm.DB) *DbFactory {
	return &DbFactory{db: db}
}
