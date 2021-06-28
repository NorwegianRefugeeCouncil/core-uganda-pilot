package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartyTypeStore struct {
	collection *mongo.Collection
}

func NewPartyTypeStore(ctx context.Context, mongoClient *mongo.Client, database string) (*PartyTypeStore, error) {
	store := &PartyTypeStore{
		collection: mongoClient.Database(database).Collection("partyTypes"),
	}

	if _, err := store.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *PartyTypeStore) Get(ctx context.Context, id string) (*PartyType, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r PartyType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *PartyTypeStore) List(ctx context.Context, listOptions PartyTypeListOptions) (*PartyTypeList, error) {

	filter := bson.M{}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*PartyType
	for {
		if !res.Next(ctx) {
			break
		}
		var r PartyType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := PartyTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *PartyTypeStore) Update(ctx context.Context, partyType *PartyType) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": partyType.ID,
	}, bson.M{
		"$set": bson.M{
			"name":      partyType.Name,
			"isBuiltIn": partyType.IsBuiltIn,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *PartyTypeStore) Create(ctx context.Context, partyType *PartyType) error {
	_, err := s.collection.InsertOne(ctx, partyType)
	if err != nil {
		return err
	}
	return nil
}
