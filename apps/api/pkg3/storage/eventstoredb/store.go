package eventstoredb

import (
	"context"
	goes "github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/client/filtering"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/position"
	"github.com/nrc-no/core/apps/api/pkg3/runtime"
	"github.com/nrc-no/core/apps/api/pkg3/storage"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"path"
	"strings"
	"time"
)

type Store struct {
	esClient    *goes.Client
	mongoClient *mongo.Client
}

func NewStore(client *goes.Client, mongoClient *mongo.Client) *Store {
	return &Store{
		esClient:    client,
		mongoClient: mongoClient,
	}
}

var _ storage.Interface = &Store{}

func (s *Store) Get(ctx context.Context, key string, out runtime.Object) error {

	keyParts := strings.Split(key, "/")
	group := keyParts[0]
	resource := keyParts[1]
	id := keyParts[2]

	collection := path.Join(group, resource)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result := s.mongoClient.Database("core").Collection(collection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return result.Err()
	}

	bytes, err := result.DecodeBytes()
	if err != nil {
		return err
	}

	raw := bytes.Lookup("currentRevision")
	if err := raw.Unmarshal(out); err != nil {
		return err
	}
	out.SetUID(id)

	return nil
}

func (s *Store) Create(ctx context.Context, key string, in, out runtime.Object) error {

	in.SetResourceVersion(1)

	response, err := s.mongoClient.Database("core").Collection("core.nrc.no/formdefinitions").InsertOne(ctx, bson.M{
		"currentRevision":   in,
		"previousRevisions": bson.A{},
		"resourceVersion":   1,
		"createdAt":         time.Now().UTC(),
		"apiVersion":        in.GetAPIVersion(),
		"group":             in.GetAPIGroup(),
		"kind":              in.GetKind(),
	})
	if err != nil {
		return err
	}

	objectID := response.InsertedID.(primitive.ObjectID)
	id := objectID.Hex()

	return s.Get(ctx, "core.nrc.no/formdefinitions/"+id, out)

}

func (s *Store) Update(ctx context.Context, key string, in, out runtime.Object) error {

	keyParts := strings.Split(key, "/")
	group := keyParts[0]
	resource := keyParts[1]
	id := keyParts[2]

	collection := path.Join(group, resource)
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.mongoClient.Database("core").Collection(collection).UpdateByID(ctx, objectID, bson.A{
		bson.M{"$set": bson.M{"previousRevisions": bson.M{"$concatArrays": bson.A{"$previousRevisions", bson.A{"$currentRevision"}}}}},
		bson.M{"$set": bson.M{"currentRevision": in}},
		bson.M{"$unset": "currentRevision.metadata.uid"},
		bson.M{"$set": bson.M{"resourceVersion": bson.M{"$add": bson.A{1, bson.M{"$size": "$previousRevisions"}}}}},
		bson.M{"$set": bson.M{"updatedAt": "$$NOW"}},
		bson.M{"$set": bson.M{"apiVersion": in.GetAPIVersion()}},
		bson.M{"$set": bson.M{"apiGroup": in.GetAPIGroup()}},
		bson.M{"$set": bson.M{"kind": in.GetKind()}},
	})
	if err != nil {
		return err
	}

	return s.Get(ctx, key, out)

}

func (s *Store) Reconcile(ctx context.Context) error {

	collection := s.mongoClient.Database("core").Collection("__meta")
	stream, err := collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	return nil
}

func (s *Store) Watch(ctx context.Context) error {

	sub, err := s.esClient.SubscribeToAllFiltered(
		ctx,
		position.StartPosition,
		false,
		filtering.NewDefaultSubscriptionFilterOptions(filtering.SubscriptionFilter{
			Prefixes: []string{"formdefinitions2-"},
			//Regex: "/.*/",
			FilterType: filtering.StreamFilter,
		}),
		func(event messages.RecordedEvent) {
			logrus.Infof("event received: %#v", event)
		}, func(p position.Position) {
			logrus.Infof("position reached: Commit: %d, Prepare: %d", p.Commit, p.Prepare)
		}, func(reason string) {
			logrus.Infof("subscription dropped: %#v", reason)
		})
	if err != nil {
		return err
	}
	if err := sub.Start(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return sub.Stop()
	}

}
