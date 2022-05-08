package entity

import (
	"context"

	"gorm.io/gorm"
)

func (s *entityPostgresStore) DeleteEntityDefinition(ctx context.Context, db *gorm.DB, entityID string) error {
	return nil
}
