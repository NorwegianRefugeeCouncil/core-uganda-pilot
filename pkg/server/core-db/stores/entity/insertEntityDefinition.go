package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresStore) InsertEntityDefinition(ctx context.Context, db *gorm.DB, entity types.EntityDefinition) (*types.EntityDefinition, error) {
	if db == nil {
		var err error
		db, err = d.db.Get()

		if err != nil {
			return nil, err
		}
	}

	ddl := d.sqlBuilder.InsertRow(
		"public",
		"entity_definition",
		[]string{
			"id",
			"name",
			"description",
			"constraint_custom",
		},
		[]any{
			entity.ID,
			entity.Name,
			entity.Description,
			entity.Constraints.Custom,
		},
	)

	result := db.Exec(ddl.Query, ddl.Args...)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &entity, nil
}
