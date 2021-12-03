package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/convert"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// DatabaseStore is the store for types.Database
type DatabaseStore interface {
	// Get a Database by ID
	Get(ctx context.Context, databaseID string) (*types.Database, error)
	// List databases
	List(ctx context.Context) (*types.DatabaseList, error)
	// Create a database
	Create(ctx context.Context, database *types.Database) (*types.Database, error)
	// Delete a database
	Delete(ctx context.Context, databaseID string) error
}

// NewDatabaseStore returns a new DatabaseStore
func NewDatabaseStore(db Factory) DatabaseStore {
	return &databaseStore{
		db:             db,
		deleteDatabase: convert.DeleteDatabaseSchemaIfExist,
		createDatabase: convert.CreateDatabase,
	}
}

// Database is the store model for types.Database
type Database struct {
	// ID corresponds to the types.Database ID
	ID string
	// Name corresponds to the types.Database Name
	Name string
	// CreatedAt corresponds to the types.Database CreatedAt
	CreatedAt time.Time
	// UpdatedAt corresponds to the types.Database UpdatedAt
	UpdatedAt time.Time
}

type deleteDatabaseFn func(db *gorm.DB, databaseId string) error
type createDatabaseFn func(db *gorm.DB, database *types.Database) error

// databaseStore is the implementation of DatabaseStore
type databaseStore struct {
	// db is the database Factory
	db Factory
	// deleteDatabase is a function for deleting the actual sql database
	// we add a variable since we want to mock this function
	deleteDatabase deleteDatabaseFn
	// createDatabase is a function for creating the actual sql database
	// we add a variable since we want to mock this function
	createDatabase createDatabaseFn
}

// Ensure that databaseStore implements DatabaseStore
var _ DatabaseStore = &databaseStore{}

// Get implements DatabaseStore.Get
func (d *databaseStore) Get(ctx context.Context, databaseID string) (*types.Database, error) {

	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "get", zap.String("database_id", databaseID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting database")
	var database Database
	if err := db.First(&database, "id = ?", databaseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Error("database not found")
			return nil, meta.NewNotFound(meta.GroupResource{
				Group:    "nrc.no",
				Resource: "databases",
			}, databaseID)
		} else {
			l.Error("failed to get database", zap.Error(err))
			return nil, meta.NewInternalServerError(err)
		}
	}

	return mapDatabaseTo(&database), nil
}

// Delete implements DatabaseStore.Delete
func (d *databaseStore) Delete(ctx context.Context, databaseID string) error {
	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "delete", zap.String("database_id", databaseID))
	if err != nil {
		return err
	}
	defer done()

	l.Debug("starting transaction")
	err = db.Transaction(func(tx *gorm.DB) error {

		l.Debug("deleting database schema")
		if err := d.deleteDatabase(tx, databaseID); err != nil {
			l.Error("failed to delete database schema", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		l.Debug("successfully deleted database schema")

		l.Debug("deleting database")
		if err := db.Delete(&Database{}, "id = ?", databaseID).Error; err != nil {
			l.Error("failed to delete database", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		l.Debug("successfully deleted database")

		l.Debug("deleting forms")
		if err := db.Delete(&Form{}, "database_id = ?", databaseID).Error; err != nil {
			l.Error("failed to delete forms", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		l.Debug("successfully deleted forms")

		l.Debug("deleting fields")
		if err := db.Delete(&Field{}, "database_id = ?", databaseID).Error; err != nil {
			l.Error("failed to delete fields", zap.Error(err))
			return meta.NewInternalServerError(err)
		}
		l.Debug("successfully deleted fields")

		return nil

	})
	l.Debug("transaction ended")

	if err != nil {
		l.Error("failed to delete detabase", zap.Error(err))
		return err
	}

	l.Debug("successfully deleted database")
	return err

}

// List implements DatabaseStore.List
func (d *databaseStore) List(ctx context.Context) (*types.DatabaseList, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "list")
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("listing databases")
	var databases []Database
	if err := db.Find(&databases).Error; err != nil {
		l.Error("failed to list databases", zap.Error(err))
		return nil, meta.NewInternalServerError(err)
	}

	l.Debug("mapping databases")
	var result = make([]*types.Database, len(databases))
	for i, database := range databases {
		result[i] = mapDatabaseTo(&database)
	}
	if result == nil {
		result = []*types.Database{}
	}

	l.Debug("successfully listed databases", zap.Int("count", len(result)))
	return &types.DatabaseList{
		Items: result,
	}, nil
}

// Create implements DatabaseStore.Create
func (d *databaseStore) Create(ctx context.Context, database *types.Database) (*types.Database, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	storeDb := mapDatabaseFrom(database)

	l.Debug("starting transaction")
	err = db.Transaction(func(tx *gorm.DB) error {

		l.Debug("storing database")
		if err := tx.Create(storeDb).Error; err != nil {
			l.Error("failed to store database", zap.Error(err))
			if IsUniqueConstraintErr(err) {
				return meta.NewAlreadyExists(meta.GroupResource{
					Group:    "core.nrc.no/v1",
					Resource: database.ID,
				}, database.ID)
			}
			return meta.NewInternalServerError(err)
		}

		l.Debug("creating database schema")
		if err := d.createDatabase(tx, database); err != nil {
			l.Error("failed to create database schema", zap.Error(err))
			tx.Rollback()
			return err
		}

		return nil

	})
	l.Debug("transaction ended")

	if err != nil {
		l.Error("failed to create database", zap.Error(err))
		return nil, err
	}

	return mapDatabaseTo(storeDb), nil
}

// mapDatabaseTo maps a store Database to a types.Database
func mapDatabaseTo(database *Database) *types.Database {
	return &types.Database{
		ID:        database.ID,
		Name:      database.Name,
		UpdatedAt: database.UpdatedAt.UTC(),
		CreatedAt: database.CreatedAt.UTC(),
	}
}

// mapDatabaseFrom maps a types.Database to a store Database
func mapDatabaseFrom(database *types.Database) *Database {
	return &Database{
		ID:   database.ID,
		Name: database.Name,
		// UpdatedAt is ignored, since the store manages this field
		// CreatedAt is ignored, since the store manages this field
	}
}
