package entity

import (
	"context"

	"github.com/nrc-no/core/pkg/common"
	store "github.com/nrc-no/core/pkg/server/core-db/stores/entity"
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
	entityStore      store.EntityStore
	transactionStore common.TransactionStore
}

func NewEntityService(entityStore store.EntityStore, transactionStore common.TransactionStore) EntityService {
	return &entityService{
		entityStore:      entityStore,
		transactionStore: transactionStore,
	}
}
