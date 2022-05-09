package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	"gorm.io/gorm"
)

func (d *entityPostgresStore) InsertRelationship(ctx context.Context, db *gorm.DB, entityRelationship types.EntityRelationship) (*types.EntityRelationship, error) {
	ddl := d.sqlBuilder.InsertRow(
		"public",
		"entity_relationship",
		[]string{
			"id",
			"cardinality",
			"source_entity_id",
			"target_entity_id",
		},
		[]any{
			entityRelationship.ID,
			entityRelationship.Cardinality,
			entityRelationship.SourceEntityID,
			entityRelationship.TargetEntityID,
		},
	)

	result := db.Exec(ddl.Query, ddl.Args...)

	if err := result.Error; err != nil {
		return nil, err
	}

	return &entityRelationship, nil
}
