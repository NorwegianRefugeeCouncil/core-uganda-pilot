package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PartyTypeStore struct {
	getCollection utils.MongoCollectionFn
}

func newPartyTypeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*PartyTypeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &PartyTypeStore{
		getCollection: utils.GetCollectionFn(database, "partyTypes", mongoClientFn),
	}

	collection, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
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
	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	res := collection.FindOne(ctx, bson.M{
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

	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filter)
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
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{
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
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, partyType)
	if err != nil {
		return err
	}
	return nil
}
