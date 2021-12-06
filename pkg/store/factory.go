package store

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Factory is the store.Factory that returns an instance of a gorm.DB
// This is useful because we can implement logic that allows us to
// renew a database connection, for example, when a connection string changes
// (credential rotation)
type Factory interface {
	Get() (*gorm.DB, error)
}

// factory is the implementation of Factory
type factory struct {
	db *gorm.DB
}

// Get implements Factory.Get
func (f factory) Get() (*gorm.DB, error) {
	return f.db, nil
}

// NewFactory returns a new instance of Factory
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

// NewMockFactory returns a mock Factory
func NewMockFactory(db *gorm.DB) Factory {
	return &factory{db: db}
}
