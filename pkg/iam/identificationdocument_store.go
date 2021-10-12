package iam

import (
	"context"
	"github.com/nrc-no/core/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IdentificationDocumentStore struct {
	getCollection utils.MongoCollectionFn
}

func newIdentificationDocumentStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*IdentificationDocumentStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &IdentificationDocumentStore{
		getCollection: utils.GetCollectionFn(database, "identificationDocuments", mongoClientFn),
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

func (s *IdentificationDocumentStore) list(ctx context.Context, listOptions IdentificationDocumentListOptions) (*IdentificationDocumentList, error) {
	filter := bson.M{}

	if len(listOptions.PartyIDs) != 0 {
		filter["partyId"] = bson.D{{Key: "$in", Value: BSONStringA(listOptions.PartyIDs)}}
	}

	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var list []*IdentificationDocument
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a IdentificationDocument
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	klose()

	return &IdentificationDocumentList{
		Items: list,
	}, nil
}

func (s *IdentificationDocumentStore) create(ctx context.Context, identificationDocument *IdentificationDocument) error {
	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, identificationDocument)
	if err != nil {
		return err
	}
	klose()
	return nil
}

func (s *IdentificationDocumentStore) get(ctx context.Context, id string) (*IdentificationDocument, error) {
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
	var a IdentificationDocument
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	klose()
	return &a, nil
}

func (s *IdentificationDocumentStore) update(ctx context.Context, identificationDocument *IdentificationDocument) error {
	collection, klose, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	_, err = collection.UpdateOne(ctx, bson.M{
		"id": identificationDocument.ID,
	}, bson.M{
		"$set": bson.M{
			"partyId":                      identificationDocument.PartyID,
			"documentNumber":               identificationDocument.DocumentNumber,
			"identificationDocumentTypeId": identificationDocument.IdentificationDocumentTypeID,
		},
	})
	if err != nil {
		return err
	}
	klose()
	return nil
}

func (s *IdentificationDocumentStore) delete(ctx context.Context, id string) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
