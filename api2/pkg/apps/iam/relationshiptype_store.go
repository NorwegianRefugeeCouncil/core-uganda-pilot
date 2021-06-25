package iam

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RelationshipTypeStore struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *RelationshipTypeStore {
	return &RelationshipTypeStore{
		collection: mongoClient.Database(database).Collection("relationshipTypes"),
	}
}

func (s *RelationshipTypeStore) Get(ctx context.Context, id string) (*RelationshipType, error) {
	res := s.collection.FindOne(ctx, bson.M{
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

	if len(listOptions.PartyType) != 0 {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"rules.firstPartyTypeId": listOptions.PartyType,
				},
				bson.M{
					"rules.secondPartyTypeId": listOptions.PartyType,
				},
			},
		}
	}

	res, err := s.collection.Find(ctx, filter)
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
	_, err := s.collection.UpdateOne(ctx, bson.M{
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
	_, err := s.collection.InsertOne(ctx, relationshipType)
	if err != nil {
		return err
	}
	return nil
}
