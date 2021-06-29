package cms

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CommentStore struct {
	collection *mongo.Collection
}

func NewCommentStore(ctx context.Context, mongoClient *mongo.Client, database string) (*CommentStore, error) {
	store := &CommentStore{
		collection: mongoClient.Database(database).Collection("comments"),
	}

	// Cases should have unique IDs
	if _, err := store.collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	// Create index on CaseID
	if _, err := store.collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys: bson.M{"caseId": 1},
		}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CommentStore) Get(ctx context.Context, id string) (*Comment, error) {
	res := s.collection.FindOne(ctx, bson.M{
		"id": id,
	})
	if res.Err() != nil {
		return nil, res.Err()
	}
	var r Comment
	if err := res.Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *CommentStore) List(ctx context.Context, options CommentListOptions) (*CommentList, error) {

	filter := bson.M{
		"caseId": options.CaseID,
	}

	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*Comment
	for {
		if !res.Next(ctx) {
			break
		}
		var r Comment
		if err := res.Decode(&r); err != nil {
			return nil, err
		}
		items = append(items, &r)
	}
	if res.Err() != nil {
		return nil, res.Err()
	}
	if items == nil {
		items = []*Comment{}
	}
	ret := CommentList{
		Items: items,
	}
	return &ret, nil
}

func (s *CommentStore) Update(ctx context.Context, id string, updateFunc func(oldComment *Comment) (*Comment, error)) (*Comment, error) {

	comment, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	newComment, err := updateFunc(comment)
	if err != nil {
		return nil, err
	}

	_, err = s.collection.ReplaceOne(ctx, bson.M{
		"id": id,
	}, newComment)
	if err != nil {
		return nil, err
	}

	return newComment, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	now := time.Now().UTC()
	comment.CreatedAt = now
	comment.UpdatedAt = now
	_, err := s.collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
