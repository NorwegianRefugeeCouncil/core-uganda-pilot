package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CaseTypeStore struct {
	getCollection utils.MongoCollectionFn
}

func NewCaseTypeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*CaseTypeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &CaseTypeStore{
		getCollection: utils.GetCollectionFn(database, "caseTypes", mongoClientFn),
	}

	collection, err := store.getCollection(ctx)
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

func (s *CaseTypeStore) Get(ctx context.Context, id string) (*CaseType, error) {
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
	var r CaseType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CaseTypeStore) List(ctx context.Context, options CaseTypeListOptions) (*CaseTypeList, error) {

	filter := bson.M{}

	if len(options.PartyTypeIDs) > 0 {
		filter["partyTypeId"] = bson.M{
			"$in": options.PartyTypeIDs,
		}
	}

	collection, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*CaseType
	for {
		if !res.Next(ctx) {
			break
		}
		var r CaseType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*CaseType{}
	}
	ret := CaseTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *CaseTypeStore) Update(ctx context.Context, caseType *CaseType) error {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": caseType.ID,
	}, bson.M{
		"$set": bson.M{
			"name":        caseType.Name,
			"partyTypeId": caseType.PartyTypeID,
			"teamId":      caseType.TeamID,
			"casTemplate": caseType.Template,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CaseTypeStore) Create(ctx context.Context, caseType *CaseType) error {
	collection, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, caseType)
	if err != nil {
		return err
	}
	return nil
}
