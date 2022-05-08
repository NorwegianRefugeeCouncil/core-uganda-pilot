package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (s *entityPostgresStore) ListEntityDefinitions(ctx context.Context, db *gorm.DB) ([]types.EntityDefinition, error) {
	return nil, nil
}
