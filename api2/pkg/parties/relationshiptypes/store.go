package relationshiptypes

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
		collection: mongoClient.Database("core").Collection("relationshipTypes"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.RelationshipType, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.RelationshipType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context) (*api.RelationshipTypeList, error) {
	res, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var items []*api.RelationshipType
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.RelationshipType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := api.RelationshipTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, relationshipType *api.RelationshipType) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": relationshipType.ID,
	}, bson.M{
		"$set": bson.M{
			"firstPartyRole":  relationshipType.FirstPartyRole,
			"secondPartyRole": relationshipType.SecondPartyRole,
			"name":            relationshipType.Name,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, relationshipType *api.RelationshipType) error {
	_, err := s.collection.InsertOne(ctx, relationshipType)
	if err != nil {
		return err
	}
	return nil
}
