package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresStore) InsertRelationship(ctx context.Context, db *gorm.DB, entityRelationship types.EntityRelationship) (*types.EntityRelationship, error) {
	return nil, nil
}
