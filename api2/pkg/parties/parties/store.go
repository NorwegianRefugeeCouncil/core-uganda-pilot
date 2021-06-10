package parties

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("parties"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.Party, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.Party
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*api.PartyList, error) {

	filter := bson.M{}

	if len(listOptions.Party) != 0 {
		filter["$or"] = bson.M{
			"firstParty":  listOptions.Party,
			"secondParty": listOptions.Party,
		}
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*api.Party
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.Party
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := api.PartyList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, party *api.Party) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": party.ID,
	}, bson.M{
		"$set": bson.M{
			"attributes": party.Attributes,
			"partyTypes": party.PartyTypes,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, party *api.Party) error {
	_, err := s.collection.InsertOne(ctx, party)
	if err != nil {
		return err
	}
	return nil
}
