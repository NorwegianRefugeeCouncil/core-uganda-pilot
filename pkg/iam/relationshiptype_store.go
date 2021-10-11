package iam

import (
	"context"
	"github.com/nrc-no/core/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RelationshipTypeStore struct {
	getCollection utils.MongoCollectionFn
}

func newRelationshipTypeStore(ctx context.Context, mongoClientFn utils.MongoClientFn, database string) (*RelationshipTypeStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &RelationshipTypeStore{
		getCollection: utils.GetCollectionFn(database, "relationshipTypes", mongoClientFn),
	}

	collection, done, err := store.getCollection(ctx)
	if err != nil {
		return nil, err
	}
	defer done()

	if _, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *RelationshipTypeStore) Get(ctx context.Context, id string) (*RelationshipType, error) {
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
	var r RelationshipType
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *RelationshipTypeStore) List(ctx context.Context, listOptions RelationshipTypeListOptions) (*RelationshipTypeList, error) {

	filter := bson.M{}

	if len(listOptions.PartyTypeID) != 0 {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"rules.firstPartyTypeId": listOptions.PartyTypeID,
				},
				bson.M{
					"rules.secondPartyTypeId": listOptions.PartyTypeID,
				},
			},
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
	var items []*RelationshipType
	for {
		if !res.Next(ctx) {
			break
		}
		var r RelationshipType
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := RelationshipTypeList{
		Items: items,
	}
	return &ret, nil
}

func (s *RelationshipTypeStore) Update(ctx context.Context, relationshipType *RelationshipType) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.UpdateOne(ctx, bson.M{
		"id": relationshipType.ID,
	}, bson.M{
		"$set": bson.M{
			"firstPartyRole":  relationshipType.FirstPartyRole,
			"secondPartyRole": relationshipType.SecondPartyRole,
			"name":            relationshipType.Name,
			"isDirectional":   relationshipType.IsDirectional,
			"rules":           relationshipType.Rules,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *RelationshipTypeStore) Create(ctx context.Context, relationshipType *RelationshipType) error {
	collection, done, err := s.getCollection(ctx)
	if err != nil {
		return err
	}
	defer done()

	_, err = collection.InsertOne(ctx, relationshipType)
	if err != nil {
		return err
	}
	return nil
}
