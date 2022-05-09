package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	uuid "github.com/satori/go.uuid"
)

func (d *entityService) Create(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	tx := d.transactionStore.Begin()

	entityId := uuid.NewV4()

	entityDefinition, err := d.entityStore.InsertEntityDefinition(ctx, tx, types.EntityDefinition{
		ID:          entityId.String(),
		Name:        entity.Name,
		Description: entity.Description,
		Constraints: entity.Constraints,
	})

	if err != nil {
		return nil, err
	}

	attributes := make([]types.Attribute, len(entity.Attributes))
	for i, attribute := range entity.Attributes {
		attributeId := uuid.NewV4()
		attribute.ID = attributeId.String()
		attribute.EntityID = entityDefinition.ID

		a, err := d.entityStore.InsertAttribute(ctx, tx, attribute)

		if err != nil {
			return nil, err
		}

		attributes[i] = *a
	}

	relationships := make([]types.EntityRelationship, len(entity.Relationships))
	for i, relationship := range entity.Relationships {
		relationshipId := uuid.NewV4()
		relationship.ID = relationshipId.String()
		relationship.SourceEntityID = entityDefinition.ID

		r, err := d.entityStore.InsertRelationship(ctx, tx, relationship)

		if err != nil {
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
