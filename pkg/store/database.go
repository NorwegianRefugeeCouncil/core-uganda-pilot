package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/convert"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type DatabaseStore interface {
	Get(databaseID string) (*types.Database, error)
	List() (*types.DatabaseList, error)
	Create(database *types.Database) (*types.Database, error)
	Delete(ctx context.Context, databaseID string) error
}

func NewDatabaseStore(db Factory) DatabaseStore {
	return &databaseStore{db: db}
}

type Database struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type databaseStore struct {
	db Factory
}

var _ DatabaseStore = &databaseStore{}

func (d *databaseStore) Get(databaseID string) (*types.Database, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	var database Database
	if err := db.First(&database, databaseID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, meta.NewNotFound(meta.GroupResource{
				Group:    "nrc.no",
				Resource: "databases",
			}, databaseID)
		} else {
			return nil, meta.NewInternalServerError(err)
		}
	}
	return mapDatabaseTo(&database), nil
}

func (d *databaseStore) Delete(ctx context.Context, databaseID string) error {

	db, err := d.db.Get()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			if err := convert.DeleteDatabaseIfExists(tx, databaseID); err != nil {
				return meta.NewInternalServerError(err)
			}
			return nil
		})

		g.Go(func() error {
			if err := db.WithContext(ctx).Delete(&Database{}, "id = ?", databaseID).Error; err != nil {
				return meta.NewInternalServerError(err)
			}
			return nil
		})

		g.Go(func() error {
			if err := db.WithContext(ctx).Delete(&Form{}, "database_id = ?", databaseID).Error; err != nil {
				return meta.NewInternalServerError(err)
			}
			return nil
		})

		g.Go(func() error {
			if err := db.WithContext(ctx).Delete(&Field{}, "database_id = ?", databaseID).Error; err != nil {
				return meta.NewInternalServerError(err)
			}
			return nil
		})

		return g.Wait()

	})

	return err

}

func (d *databaseStore) List() (*types.DatabaseList, error) {
	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	var databases []Database
	if err := db.Find(&databases).Error; err != nil {
		return nil, meta.NewInternalServerError(err)
	}
	var result = make([]*types.Database, len(databases))
	for i, database := range databases {
		result[i] = mapDatabaseTo(&database)
	}
	return &types.DatabaseList{
		Items: result,
	}, nil
}

func (d *databaseStore) Create(database *types.Database) (*types.Database, error) {

	db, err := d.db.Get()
	if err != nil {
		return nil, err
	}

	database.ID = uuid.NewV4().String()

	storeDb := mapDatabaseFrom(database)
	storeDb.CreatedAt = time.Now()

	err = db.Transaction(func(tx *gorm.DB) error {

		if err := db.Create(storeDb).Error; err != nil {
			return meta.NewInternalServerError(err)
		}

		sqlDB, err := tx.DB()
		if err != nil {
			return err
		}

		if err := convert.CreateDatabase(sqlDB, database); err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return mapDatabaseTo(storeDb), nil
}

func mapDatabaseTo(database *Database) *types.Database {
	return &types.Database{
		ID:   database.ID,
		Name: database.Name,
	}
}

func mapDatabaseFrom(database *types.Database) *Database {
	return &Database{
		ID:   database.ID,
		Name: database.Name,
	}
}
