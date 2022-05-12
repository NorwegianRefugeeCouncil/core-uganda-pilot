package foo

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"github.com/nrc-no/core/pkg/sqlmanager"
	"github.com/nrc-no/core/pkg/store"
)

type EntityModel interface {
	Create(ctx context.Context, entity types.Entity) (*types.Entity, error)
}

type entityPostgresModel struct {
	db         store.Factory
	sqlBuilder sqlmanager.SQLBuilder
}

func NewEntityPostgresModel(db store.Factory) EntityModel {
	return &entityPostgresModel{
		db:         db,
		sqlBuilder: sqlmanager.NewSQLBuilder(),
	}
}
