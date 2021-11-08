package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Factory interface {
	Get() (*gorm.DB, error)
}

type factory struct {
	db *gorm.DB
}

func (f factory) Get() (*gorm.DB, error) {
	return f.db, nil
}

func NewFactory(dsn string) (Factory, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}))
	if err != nil {
		return nil, err
	}
	return &factory{
		db: db,
	}, nil
}

func NewMockFactory(db *gorm.DB) Factory {
	return &factory{db: db}
}
