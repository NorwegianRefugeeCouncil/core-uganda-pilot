package parties

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	Collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *Store {
	return &Store{
		Collection: mongoClient.Database(database).Collection("parties"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*Party, error) {
	res := s.Collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Party
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*PartyList, error) {
	filter := bson.M{}

	if len(listOptions.PartyTypeID) > 0 {
		filter["partyTypes"] = listOptions.PartyTypeID
	}

	res, err := s.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Party
	for {
		if !res.Next(ctx) {
			break
		}
		var r Party
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Party{}
	}
	ret := PartyList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, party *Party) error {
	_, err := s.Collection.UpdateOne(ctx, bson.M{
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

func (s *Store) Create(ctx context.Context, party *Party) error {
	_, err := s.Collection.InsertOne(ctx, party)
	if err != nil {
		return err
	}
	return nil
}
