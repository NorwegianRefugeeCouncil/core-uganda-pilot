package relationships

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(mongoClient *mongo.Client, database string) *Store {
	return &Store{
		collection: mongoClient.Database(database).Collection("relationships"),
	}
}

func (s *Store) Get(ctx context.Context, id string) (*Relationship, error) {
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

func (s *Store) List(ctx context.Context, listOptions ListOptions) (*RelationshipList, error) {

	var filterExpressions []bson.M

	if len(listOptions.RelationshipTypeID) != 0 {
		filterExpressions = append(filterExpressions, bson.M{
			"relationshipTypeId": listOptions.RelationshipTypeID,
		})
	}

	if len(listOptions.EitherParty) != 0 {
		filterExpressions = append(filterExpressions, bson.M{
			"$or": bson.A{
				bson.M{"firstParty": listOptions.EitherParty},
				bson.M{"secondParty": listOptions.EitherParty},
			},
		})
	} else {
		if len(listOptions.FirstPartyId) != 0 {
			filterExpressions = append(filterExpressions, bson.M{"firstParty": listOptions.FirstPartyId})
		}
		if len(listOptions.SecondParty) != 0 {
			filterExpressions = append(filterExpressions, bson.M{"secondParty": listOptions.SecondParty})

		}
	}

	filter := bson.M{}

	if len(filterExpressions) > 0 {
		expressions := bson.A{}
		for _, filterExpression := range filterExpressions {
			expressions = append(expressions, filterExpression)
		}
		filter["$and"] = expressions
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

func (s *Store) Update(ctx context.Context, relationship *Relationship) error {
	_, err := s.collection.UpdateOne(ctx, bson.M{
		"id": relationship.ID,
	}, bson.M{
		"$set": bson.M{
			"startOfRelationship": relationship.StartOfRelationship,
			"endOfRelationship":   relationship.EndOfRelationship,
			"firstParty":          relationship.FirstParty,
			"secondParty":         relationship.SecondParty,
			"relationshipTypeId":  relationship.RelationshipTypeID,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Create(ctx context.Context, relationship *Relationship) error {
	_, err := s.collection.InsertOne(ctx, relationship)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
