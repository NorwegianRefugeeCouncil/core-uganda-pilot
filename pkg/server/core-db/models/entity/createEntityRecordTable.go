package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresModel) CreateEntityRecordTable(ctx context.Context, db *gorm.DB, entity types.Entity) error {
	return nil
}
