package response

import (
	"context"
	"github.com/nrc-no/core-kafka/pkg/response/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	collection *mongo.Collection
}

func (s *Store) FindPendingResponseRequestFor(
	ctx context.Context,
	subjectType,
	subject string,
) (*api.ResponseRequest, error) {

	result := s.collection.FindOne(ctx, bson.M{
		"subjecttype": subjectType,
		"subject":     subject,
	}, options.FindOne().SetSort(bson.D{{"createdAt", -1}}))

	if result.Err() != nil {
		return nil, result.Err()
	}

	ret := api.ResponseRequest{}

}
