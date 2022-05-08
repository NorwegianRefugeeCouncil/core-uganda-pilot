package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (s *entityPostgresStore) UpdateEntityDefinition(ctx context.Context, db *gorm.DB, entity types.EntityDefinition) (*types.EntityDefinition, error) {
	return nil, nil
}
