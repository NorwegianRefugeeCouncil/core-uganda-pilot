package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttributeStore struct {
	collection *mongo.Collection
}

func NewAttributeStore(ctx context.Context, mongoClient *mongo.Client, database string) (*AttributeStore, error) {
	store := &AttributeStore{
		collection: mongoClient.Database(database).Collection("attributes"),
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

func (s *AttributeStore) List(ctx context.Context, listOptions AttributeListOptions) (*AttributeList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypeIDs) > 0 {
		filter["partyTypeIds"] = bson.M{
			"$all": listOptions.PartyTypeIDs,
		}
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*Attribute
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a Attribute
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &AttributeList{
		Items: list,
	}, nil

}

func (s *AttributeStore) Create(ctx context.Context, attribute *Attribute) error {
	_, err := s.collection.InsertOne(ctx, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *AttributeStore) Get(ctx context.Context, id string) (*Attribute, error) {
	result := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var a Attribute
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AttributeStore) Update(ctx context.Context, attribute *Attribute) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": attribute.ID,
	}, bson.M{
		"$set": bson.M{
			"name":         attribute.Name,
			"translations": attribute.Translations,
			"isPii":        attribute.IsPersonallyIdentifiableInfo,
			"partyTypeIds": attribute.PartyTypeIDs,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
