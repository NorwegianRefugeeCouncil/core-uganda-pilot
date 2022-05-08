package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (s *entityPostgresStore) UpdateAttribute(ctx context.Context, db *gorm.DB, attribute types.Attribute) (*types.Attribute, error) {
	return nil, nil
}
