package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (s *entityPostgresStore) UpdateRelationship(ctx context.Context, db *gorm.DB, relationship types.EntityRelationship) (*types.EntityRelationship, error) {
	return nil, nil
}
