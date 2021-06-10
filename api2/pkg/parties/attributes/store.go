package attributes

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
		collection: mongoClient.Database("core").Collection("attributes"),
	}
}

func (s *Store) List(ctx context.Context) (*api.AttributeList, error) {

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var list []*api.Attribute
	for {
		if !cursor.Next(ctx) {
			break
		}
		var a api.Attribute
		if err := cursor.Decode(&a); err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &api.AttributeList{
		Items: list,
	}, nil

}

func (s *Store) Create(ctx context.Context, attribute *api.Attribute) error {
	_, err := s.collection.InsertOne(ctx, attribute)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (*api.Attribute, error) {
	result := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var a api.Attribute
	if err := result.Decode(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Store) Update(ctx context.Context, attribute *api.Attribute) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": attribute.ID,
	}, bson.M{
		"$set": bson.M{
			"name":                         attribute.Name,
			"translations":                 attribute.Translations,
			"valueType":                    attribute.ValueType,
			"isPersonallyIdentifiableInfo": attribute.IsPersonallyIdentifiableInfo,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
