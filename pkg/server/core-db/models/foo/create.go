package foo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/nrc-no/core/pkg/server/core-db/types"
	"github.com/nrc-no/core/pkg/sqlmanager"
	"gorm.io/gorm"
)

func (d *entityPostgresModel) Create(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	tx := d.db.Begin()

	entityID := uuid.NewV4()

	entityDefinition, err := insertEntityDefinition(ctx, tx, types.EntityDefinition{
		ID:          entityID.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Constraints: entity.Constraints,
	})

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	attributes := make([]types.Attribute, len(entity.Attributes))
	for i, attribute := range entity.Attributes {
		attributeId := uuid.NewV4()
		attribute.ID = attributeId.String()
		attribute.EntityID = entityID

		a, err := insertAttribute(ctx, tx, attribute, d.sqlBuilder)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		attributes[i] = *a
	}

	relationships := make([]types.EntityRelationship, len(entity.Relationships))
	for i, relationship := range entity.Relationships {
		relationshipId := uuid.NewV4()
		relationship.ID = relationshipId.String()
		relationship.SourceEntityID = entityID

		r, err := insertRelationship(ctx, tx, relationship, d.sqlBuilder)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		relationships[i] = *r
	}

	createdEntity := types.Entity{
		EntityDefinition: *entityDefinition,
		Attributes:       attributes,
		Relationships:    relationships,
	}

	return &createdEntity, nil
}

func insertEntityDefinition(ctx context.Context, db *gorm.DB, entity types.EntityDefinition, sqlBuilder sqlmanager.SQLBuilder) (*types.EntityDefinition, error) {
	ddl := sqlBuilder.InsertRow(
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

func insertAttribute(ctx context.Context, db *gorm.DB, attribute types.Attribute, sqlBuilder sqlmanager.SQLBuilder) (*types.Attribute, error) {
	ddl := sqlBuilder.InsertRow(
		"public",
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

func insertRelationship(ctx context.Context, db *gorm.DB, entityRelationship types.EntityRelationship, sqlBuilder sqlmanager.SQLBuilder) (*types.EntityRelationship, error) {
	ddl := sqlBuilder.InsertRow(
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
