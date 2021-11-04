package store

import (
	"context"
	"errors"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/sql/convert"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type DatabaseStore interface {
	Get(ctx context.Context, databaseID string) (*types.Database, error)
	List(ctx context.Context) (*types.DatabaseList, error)
	Create(ctx context.Context, database *types.Database) (*types.Database, error)
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

func (d *databaseStore) Get(ctx context.Context, databaseID string) (*types.Database, error) {

	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "get", zap.String("database_id", databaseID))
	if err != nil {
		return nil, err
	}
	defer done()

	l.Debug("getting database")
	var database Database
	if err := db.First(&database, databaseID).Error; err != nil {
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

func (d *databaseStore) Delete(ctx context.Context, databaseID string) error {
	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "delete", zap.String("database_id", databaseID))
	if err != nil {
		return err
	}
	defer done()

	l.Debug("starting transaction")
	err = db.Transaction(func(tx *gorm.DB) error {

		g, ctx := errgroup.WithContext(ctx)

		g.Go(func() error {
			l.Debug("deleting database schema")
			if err := convert.DeleteDatabaseSchemaIfExist(tx, databaseID); err != nil {
				l.Error("failed to delete database schema", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("successfully deleted database schema")
			return nil
		})

		g.Go(func() error {
			l.Debug("deleting database")
			if err := db.WithContext(ctx).Delete(&Database{}, "id = ?", databaseID).Error; err != nil {
				l.Error("failed to delete database", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("successfully deleted database")
			return nil
		})

		g.Go(func() error {
			l.Debug("deleting forms")
			if err := db.WithContext(ctx).Delete(&Form{}, "database_id = ?", databaseID).Error; err != nil {
				l.Error("failed to delete forms", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("successfully deleted forms")
			return nil
		})

		g.Go(func() error {
			l.Debug("deleting fields")
			if err := db.WithContext(ctx).Delete(&Field{}, "database_id = ?", databaseID).Error; err != nil {
				l.Error("failed to delete fields", zap.Error(err))
				return meta.NewInternalServerError(err)
			}
			l.Debug("successfully deleted fields")
			return nil
		})

		if err := g.Wait(); err != nil {
			l.Error("failed to delete database", zap.Error(err))
			return err
		}

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

func (d *databaseStore) Create(ctx context.Context, database *types.Database) (*types.Database, error) {
	ctx, db, l, done, err := actionContext(ctx, d.db, "database", "create")
	if err != nil {
		return nil, err
	}
	defer done()

	database.ID = uuid.NewV4().String()
	storeDb := mapDatabaseFrom(database)
	storeDb.CreatedAt = time.Now()

	l.Debug("starting tansaction")
	err = db.Transaction(func(tx *gorm.DB) error {

		l.Debug("storing database")
		if err := db.Create(storeDb).Error; err != nil {
			l.Error("failed to store database", zap.Error(err))
			return meta.NewInternalServerError(err)
		}

		l.Debug("creating database schema")
		if err := convert.CreateDatabase(tx, database); err != nil {
			l.Error("failed to create database schema", zap.Error(err))
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
