package partytypeschemas

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/parties/partytypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *Store {
	return &Store{
		collection: mongoClient.Database(database).Collection("partyTypeSchemas"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*partytypes.PartyTypeSchema, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r partytypes.PartyTypeSchema
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *Store) List(ctx context.Context) (*partytypes.PartyTypeSchemaList, error) {

	filter := bson.M{}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*partytypes.PartyTypeSchema
	for {
		if !res.Next(ctx) {
			break
		}
		var r partytypes.PartyTypeSchema
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	ret := partytypes.PartyTypeSchemaList{
		Items: items,
	}
	return &ret, nil
}

func (s *Store) Update(ctx context.Context, partyTypeSchema *partytypes.PartyTypeSchema) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": partyTypeSchema.ID,
	}, bson.M{
		"$set": bson.M{
			"name":  partyTypeSchema.Name,
			"nodes": partyTypeSchema.Nodes,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, partyTypeSchema *partytypes.PartyTypeSchema) error {
	_, err := s.collection.InsertOne(ctx, partyTypeSchema)
	if err != nil {
		return err
	}
	return nil
}
