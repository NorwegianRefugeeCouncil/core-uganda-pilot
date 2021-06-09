package store

import (
	"context"
	v1 "github.com/nrc-no/core-kafka/pkg/collection/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Topic struct {
	mongoClient *mongo.Client
}

func NewTopicStore(mongoClient *mongo.Client) *Topic {
	return &Topic{
		mongoClient: mongoClient,
	}
}

func (t *Topic) GetTopic(ctx context.Context, topic v1.Topic) (*v1.TopicDescription, error) {
	result := t.mongoClient.Database("core").Collection("topics").FindOne(ctx, bson.M{
		"topic": topic,
	})
	if result.Err() != nil {
		return nil, result.Err()
	}
	var ret v1.TopicDescription
	if err := result.Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (t *Topic) CreateTopic(ctx context.Context, topic *v1.TopicDescription) error {
	_, err := t.mongoClient.Database("core").Collection("topics").InsertOne(ctx, topic)
	if err != nil {
		return err
	}
	return nil
}

func (t *Topic) ReplaceTopic(ctx context.Context, topic *v1.TopicDescription) error {
	_, err := t.mongoClient.Database("core").Collection("topics").ReplaceOne(
		ctx,
		bson.M{
			"topic": topic.Topic,
		},
		topic)
	if err != nil {
		return err
	}
	return nil
}
