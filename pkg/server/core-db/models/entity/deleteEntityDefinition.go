package entity

import (
	"context"

	"gorm.io/gorm"
)

func (s *entityPostgresModel) DeleteEntityDefinition(ctx context.Context, db *gorm.DB, entityID string) error {
	return nil
}
