package relationshipparties

import (
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
)

type PartiesStore struct {
	store *parties.Store
}

func NewStore(partiesStore *parties.Store) *PartiesStore {
	return &PartiesStore{
		store: partiesStore,
	}
}
