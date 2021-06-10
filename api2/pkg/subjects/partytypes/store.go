package partytypes

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/subjects/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("partyTypes"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.PartyType, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.PartyType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context) (*api.PartyTypeList, error) {

	filter := bson.M{}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*api.PartyType
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.PartyType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := api.PartyTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, partyType *api.PartyType) error {
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

func (s *Store) Create(ctx context.Context, partyType *api.PartyType) error {
	_, err := s.collection.InsertOne(ctx, partyType)
	if err != nil {
		return err
	}
	return nil
}
