package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RelationshipStore struct {
	collection *mongo.Collection
}

func NewRelationshipStore(ctx context.Context, mongoClient *mongo.Client, database string) (*RelationshipStore, error) {
	store := &RelationshipStore{
		collection: mongoClient.Database(database).Collection("relationships"),
	}

	if _, err := store.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	return store, nil

}

func (s *RelationshipStore) Get(ctx context.Context, id string) (*Relationship, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Relationship
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RelationshipStore) List(ctx context.Context, listOptions RelationshipListOptions) (*RelationshipList, error) {

	filter := bson.M{}

	if len(listOptions.RelationshipTypeID) != 0 {
		filter["relationshipTypeId"] = listOptions.RelationshipTypeID
	}

	if len(listOptions.EitherPartyID) != 0 {
		filter["$or"] = bson.A{
			bson.M{"firstParty": listOptions.EitherPartyID},
			bson.M{"secondParty": listOptions.EitherPartyID},
		}
	} else {
		if len(listOptions.FirstPartyID) != 0 {
			filter["firstParty"] = listOptions.FirstPartyID
		}
		if len(listOptions.SecondPartyID) != 0 {
			filter["secondParty"] = listOptions.SecondPartyID

		}
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Relationship
	for {
		if !res.Next(ctx) {
			break
		}
		var r Relationship
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := RelationshipList{
		Items: items,
	}
	return &ret, nil
}

func (s *RelationshipStore) Update(ctx context.Context, relationship *Relationship) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": relationship.ID,
	}, bson.M{
		"$set": bson.M{
			"firstParty":         relationship.FirstPartyID,
			"secondParty":        relationship.SecondPartyID,
			"relationshipTypeId": relationship.RelationshipTypeID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *RelationshipStore) Create(ctx context.Context, relationship *Relationship) error {
	_, err := s.collection.InsertOne(ctx, relationship)
	if err != nil {
		return err
	}
	return nil
}

func (s *RelationshipStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
