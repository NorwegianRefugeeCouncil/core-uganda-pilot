package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CaseStore struct {
	getCollection func() (*mongo.Collection, error)
}

func NewCaseStore(ctx context.Context, mongoClientSrc storage.MongoClientSrc, database string) (*CaseStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &CaseStore{
		getCollection: storage.GetCollectionFn(mongoClientSrc, database, "cases"),
	}

	collection, err := store.getCollection()
	if err != nil {
		return nil, err
	}

	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CaseStore) create(ctx context.Context, kase *Case) error {
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, kase)
	if err != nil {
		return err
	}
	return nil
}

func (s *CaseStore) get(ctx context.Context, id string) (*Case, error) {
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res := collection.FindOne(ctx, bson.M{
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

func (s *CaseStore) list(ctx context.Context, listOptions CaseListOptions) (*CaseList, error) {
	filter := bson.M{}

	if len(listOptions.PartyIDs) > 0 {
		filter["partyId"] = bson.M{"$in": listOptions.PartyIDs}
	}
	if len(listOptions.CaseTypeIDs) > 0 {
		filter["caseTypeId"] = bson.M{"$in": listOptions.CaseTypeIDs}
	}
	if len(listOptions.TeamIDs) > 0 {
		filter["teamId"] = bson.M{"$in": listOptions.TeamIDs}
	}
	if len(listOptions.ParentID) > 0 {
		filter["parentId"] = listOptions.ParentID
	}
	if listOptions.Done != nil {
		filter["done"] = *listOptions.Done
	}

	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filter)
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

func (s *CaseStore) update(ctx context.Context, kase *Case) error {
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": kase.ID,
	}, bson.M{
		"$set": bson.M{
			"formData": kase.FormData,
			"done":     kase.Done,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
