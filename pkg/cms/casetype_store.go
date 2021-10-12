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

	collection, done, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CaseTypeStore) get(ctx context.Context, id string) (*CaseType, error) {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

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

func (s *CaseTypeStore) list(ctx context.Context, options CaseTypeListOptions) (*CaseTypeList, error) {

	filter := bson.M{}

	if len(options.PartyTypeIDs) > 0 {
		filter["partyTypeId"] = bson.M{
			"$in": options.PartyTypeIDs,
		}
	}

	if len(options.TeamIDs) > 0 {
		filter["teamId"] = bson.M{
			"$in": options.TeamIDs,
		}
	}

	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

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

func (s *CaseTypeStore) update(ctx context.Context, caseType *CaseType) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": caseType.ID,
	}, bson.M{
		"$set": bson.M{
			"name":           caseType.Name,
			"partyTypeId":    caseType.PartyTypeID,
			"teamId":         caseType.TeamID,
			"form":           caseType.Form,
			"intakeCaseType": caseType.IntakeCaseType,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CaseTypeStore) create(ctx context.Context, caseType *CaseType) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.InsertOne(ctx, caseType)
	if err != nil {
		return err
	}
	return nil
}
