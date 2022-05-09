package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/server/core-db/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (d *entityService) Create(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	tx := d.transactionManager.Begin()

	entityId := uuid.NewV4()

	entityDefinition, err := d.entityModel.InsertEntityDefinition(ctx, tx, types.EntityDefinition{
		ID:          entityId.String(),
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
		a, err := d.createAttribute(ctx, tx, attribute, entityDefinition.ID)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		attributes[i] = *a
	}

	relationships := make([]types.EntityRelationship, len(entity.Relationships))
	for i, relationship := range entity.Relationships {
		r, err := d.createEntityRelationship(ctx, tx, relationship, entityDefinition.ID)

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

func (d *entityService) createAttribute(ctx context.Context, tx *gorm.DB, attribute types.Attribute, entityID string) (*types.Attribute, error) {
	attributeId := uuid.NewV4()
	attribute.ID = attributeId.String()
	attribute.EntityID = entityID

	a, err := d.entityModel.InsertAttribute(ctx, tx, attribute)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (d *entityService) createEntityRelationship(ctx context.Context, tx *gorm.DB, relationship types.EntityRelationship, entityID string) (*types.EntityRelationship, error) {
	relationshipId := uuid.NewV4()
	relationship.ID = relationshipId.String()
	relationship.SourceEntityID = entityID

	r, err := d.entityModel.InsertRelationship(ctx, tx, relationship)

	if err != nil {
		return nil, err
	}

	return r, nil
}
