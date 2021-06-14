package casetypes

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/cases/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("caseTypes"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.CaseType, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.CaseType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context) (*api.CaseTypeList, error) {

	filter := bson.M{}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*api.CaseType
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.CaseType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*api.CaseType{}
	}
	ret := api.CaseTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, caseType *api.CaseType) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": caseType.ID,
	}, bson.M{
		"$set": bson.M{
			"name":        caseType.Name,
			"partyTypeId": caseType.PartyTypeID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, caseType *api.CaseType) error {
	_, err := s.collection.InsertOne(ctx, caseType)
	if err != nil {
		return err
	}
	return nil
}
