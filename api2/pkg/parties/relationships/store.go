package relationships

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("relationships"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*api.Relationship, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r api.Relationship
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*api.RelationshipList, error) {

	filter := bson.M{}

	if len(listOptions.Party) != 0 {
		filter["$or"] = bson.M{
			"firstParty":  listOptions.Party,
			"secondParty": listOptions.Party,
		}
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*api.Relationship
	for {
		if !res.Next(ctx) {
			break
		}
		var r api.Relationship
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := api.RelationshipList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, relationship *api.Relationship) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": relationship.ID,
	}, bson.M{
		"$set": bson.M{
			"startOfRelationship": relationship.StartOfRelationship,
			"endOfRelationship":   relationship.EndOfRelationship,
			"firstParty":          relationship.FirstParty,
			"secondParty":         relationship.SecondParty,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, relationship *api.Relationship) error {
	_, err := s.collection.InsertOne(ctx, relationship)
	if err != nil {
		return err
	}
	return nil
}
