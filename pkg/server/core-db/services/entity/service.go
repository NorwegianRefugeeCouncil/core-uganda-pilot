package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/common"
	model "github.com/nrc-no/core/pkg/server/core-db/models/entity"
	"github.com/nrc-no/core/pkg/server/core-db/types"
)

type EntityService interface {
	Create(ctx context.Context, entity types.Entity) (*types.Entity, error)
	Get(ctx context.Context, entityID string) (*types.Entity, error)
	List(ctx context.Context) ([]types.EntityDefinition, error)
	Update(ctx context.Context, entity types.Entity) (*types.Entity, error)
	Delete(ctx context.Context, entityID string) error
}

type entityService struct {
	entityModel        model.EntityModel
	transactionManager common.TransactionManager
}

func NewEntityService(entityModel model.EntityModel, transactionManager common.TransactionManager) EntityService {
	return &entityService{
		entityModel:        entityModel,
		transactionManager: transactionManager,
	}
}
