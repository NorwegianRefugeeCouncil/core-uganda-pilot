package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"github.com/nrc-no/core/pkg/sqlmanager"
	"github.com/nrc-no/core/pkg/store"
	"gorm.io/gorm"
)

type EntityStore interface {
	ListEntityDefinitions(ctx context.Context, db *gorm.DB) ([]types.EntityDefinition, error)

	GetEntity(ctx context.Context, db *gorm.DB, entityID string) (*types.Entity, error)

	InsertEntityDefinition(ctx context.Context, db *gorm.DB, entity types.EntityDefinition) (*types.EntityDefinition, error)
	InsertAttribute(ctx context.Context, db *gorm.DB, attribute types.Attribute) (*types.Attribute, error)
	InsertRelationship(ctx context.Context, db *gorm.DB, relationship types.EntityRelationship) (*types.EntityRelationship, error)
	CreateEntityRecordTable(ctx context.Context, db *gorm.DB, entity types.Entity) error

	UpdateEntityDefinition(ctx context.Context, db *gorm.DB, entity types.EntityDefinition) (*types.EntityDefinition, error)
	UpdateAttribute(ctx context.Context, db *gorm.DB, attribute types.Attribute) (*types.Attribute, error)
	UpdateRelationship(ctx context.Context, db *gorm.DB, relationship types.EntityRelationship) (*types.EntityRelationship, error)

	DeleteEntityDefinition(ctx context.Context, db *gorm.DB, entityID string) error
	DropEntityRecordTable(ctx context.Context, db *gorm.DB, entityID string) error
}

type entityPostgresStore struct {
	db         store.Factory
	sqlBuilder sqlmanager.SQLBuilder
}

func NewEntityPostgresStore(db store.Factory) EntityStore {
	return &entityPostgresStore{
		db:         db,
		sqlBuilder: sqlmanager.NewSQLBuilder(),
	}
}
