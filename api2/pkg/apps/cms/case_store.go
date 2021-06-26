package cms

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CaseStore struct {
	collection *mongo.Collection
}

func NewCaseStore(ctx context.Context, mongoClient *mongo.Client, database string) (*CaseStore, error) {
	store := &CaseStore{
		collection: mongoClient.Database(database).Collection("cases"),
	}

	if _, err := store.collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CaseStore) Get(ctx context.Context, id string) (*Case, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Case
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CaseStore) List(ctx context.Context, listOptions CaseListOptions) (*CaseList, error) {

	filter := bson.M{}

	if len(listOptions.PartyID) > 0 {
		filter["partyId"] = listOptions.PartyID
	}
	if len(listOptions.CaseTypeID) > 0 {
		filter["caseTypeId"] = listOptions.CaseTypeID
	}
	if len(listOptions.ParentID) > 0 {
		filter["parentId"] = listOptions.ParentID
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Case
	for {
		if !res.Next(ctx) {
			break
		}
		var r Case
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Case{}
	}
	ret := CaseList{
		Items: items,
	}
	return &ret, nil
}

func (s *CaseStore) Update(ctx context.Context, kase *Case) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": kase.ID,
	}, bson.M{
		"$set": bson.M{
			"description": kase.Description,
			"done":        kase.Done,
			"parentId":    kase.ParentID,
			"teamId":      kase.TeamID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CaseStore) Create(ctx context.Context, kase *Case) error {
	_, err := s.collection.InsertOne(ctx, kase)
	if err != nil {
		return err
	}
	return nil
}
