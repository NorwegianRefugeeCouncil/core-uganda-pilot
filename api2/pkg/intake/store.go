package intake

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(
	mongoClient *mongo.Client,
) *Store {
	return &Store{
		collection: mongoClient.Database("core").Collection("submissions"),
	}
}

func (s *Store) StoreSubmission(ctx context.Context, submission *Submission) error {
	_, err := s.collection.InsertOne(ctx, submission)
	if err != nil {
		return err
	}
	return nil
}
