package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartyTypeStore struct {
	collection *mongo.Collection
}

func NewPartyTypeStore(mongoClient *mongo.Client, database string) *PartyTypeStore {
	return &PartyTypeStore{
		collection: mongoClient.Database(database).Collection("partyTypes"),
	}
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

func (s *PartyTypeStore) List(ctx context.Context) (*PartyTypeList, error) {

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
