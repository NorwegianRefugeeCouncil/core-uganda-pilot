package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IdentificationDocumentTypeStore struct {
	getCollection utils.MongoCollectionFn
}

func newIdentificationDocumentTypeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*IdentificationDocumentTypeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &IdentificationDocumentTypeStore{
		getCollection: utils.GetCollectionFn(database, IdentificationDocumentTypesCollection, mongoClientFn),
	}

	collection, klose, err := store.getCollection(ctx)
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

	klose()

	return store, nil
}

func (s *IdentificationDocumentTypeStore) list(ctx context.Context, listOptions IdentificationDocumentTypeListOptions) (*IdentificationDocumentTypeList, error) {
	filter := bson.M{}

	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*IdentificationDocumentType
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a IdentificationDocumentType
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	klose()

	return &IdentificationDocumentTypeList{
		Items: list,
	}, nil
}

func (s *IdentificationDocumentTypeStore) create(ctx context.Context, identificationDocumentType *IdentificationDocumentType) error {
	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, identificationDocumentType)
	if err != nil {
		return err
	}
	klose()
	return nil
}

func (s *IdentificationDocumentTypeStore) get(ctx context.Context, id string) (*IdentificationDocumentType, error) {
	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var a IdentificationDocumentType
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	klose()
	return &a, nil
}

func (s *IdentificationDocumentTypeStore) update(ctx context.Context, identificationDocumentType *IdentificationDocumentType) error {
	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{
		"id": identificationDocumentType.ID,
	}, bson.M{
		"$set": bson.M{
			"name": identificationDocumentType.Name,
		},
	})
	if err != nil {
		return err
	}
	klose()
	return nil
}
