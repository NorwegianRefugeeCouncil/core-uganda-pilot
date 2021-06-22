package casetypes

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *Store {
	return &Store{
		collection: mongoClient.Database(database).Collection("caseTypes"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*CaseType, error) {
	res := s.collection.FindOne(ctx, bson.M{
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

func (s *Store) List(ctx context.Context, options ListOptions) (*CaseTypeList, error) {

	filter := bson.M{}

	if len(options.PartyTypeIDs) > 0 {
		filter["partyTypeId"] = bson.M{
			"$in": options.PartyTypeIDs,
		}
	}

	res, err := s.collection.Find(ctx, filter)
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

func (s *Store) Update(ctx context.Context, caseType *CaseType) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": caseType.ID,
	}, bson.M{
		"$set": bson.M{
			"name":        caseType.Name,
			"partyTypeId": caseType.PartyTypeID,
			"teamId":      caseType.TeamID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, caseType *CaseType) error {
	_, err := s.collection.InsertOne(ctx, caseType)
	if err != nil {
		return err
	}
	return nil
}
