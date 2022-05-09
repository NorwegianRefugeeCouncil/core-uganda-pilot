package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresModel) GetEntity(ctx context.Context, db *gorm.DB, entityID string) (*types.Entity, error) {
	return nil, nil
}
