package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresStore) InsertAttribute(ctx context.Context, db *gorm.DB, attribute types.Attribute) (*types.Attribute, error) {
	ddl := d.sqlBuilder.InsertRow(
		"",
		"entity_attribute",
		[]string{
			"id",
			"name",
			"list",
			"type",
			"entity_id",
			"constraint_required",
			"constraint_unique",
			"constraint_min",
			"constraint_max",
			"constraint_pattern",
			"constraint_enum",
			"constraint_custom",
		},
		[]any{
			attribute.ID,
			attribute.Name,
			attribute.List,
			attribute.Type,
			attribute.EntityID,
			attribute.Constraints.Required,
			attribute.Constraints.Unique,
			attribute.Constraints.Min,
			attribute.Constraints.Max,
			attribute.Constraints.Pattern,
			attribute.Constraints.Enum,
			attribute.Constraints.Custom,
		},
	)

	result := db.Exec(ddl.Query, ddl.Args...)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &attribute, nil
}
