package relationshipparties

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/parties"
	"go.mongodb.org/mongo-driver/bson"
)

type PartiesStore struct {
	store *parties.Store
}

func NewStore(partiesStore *parties.Store) *PartiesStore {
	return &PartiesStore{
		store: partiesStore,
	}
}

func (s *PartiesStore) FilteredList(ctx context.Context, filterOptions PickPartyOptions) (*parties.PartyList, error) {
	filter := bson.M{}
	res, err := s.store.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*parties.Party
	for {
		if !res.Next(ctx) {
			break
		}
		var r parties.Party
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*parties.Party{}
	}
	ret := parties.PartyList{
		Items: items,
	}
	return &ret, nil
}