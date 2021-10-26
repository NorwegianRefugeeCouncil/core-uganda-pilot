package cms

import (
	"context"
	"github.com/nrc-no/core/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CommentStore struct {
	getCollection func() (*mongo.Collection, error)
}

func NewCommentStore(ctx context.Context, mongoClientSrc storage.MongoClientSrc, database string) (*CommentStore, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	store := &CommentStore{
		getCollection: storage.GetCollectionFn(mongoClientSrc, database, "comments"),
	}

	collection, err := store.getCollection()
	if err != nil {
		return nil, err
	}

	// Cases should have unique IDs
	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		}); err != nil {
		return nil, err
	}

	// Create index on CaseID
	if _, err := collection.Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys: bson.M{"caseId": 1},
		}); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *CommentStore) Get(ctx context.Context, id string) (*Comment, error) {
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res := collection.FindOne(ctx, bson.M{
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

	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	res, err := collection.Find(ctx, filter)
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

	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	_, err = collection.ReplaceOne(ctx, bson.M{
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

	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, comment)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) Delete(ctx context.Context, id string) error {
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{
		"id": id,
	})
	if err != nil {
		return err
	}
	return nil
}
