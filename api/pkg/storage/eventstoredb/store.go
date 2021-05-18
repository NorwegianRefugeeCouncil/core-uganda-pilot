package eventstoredb

import (
	"context"
	goes "github.com/EventStore/EventStore-Client-Go/client"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"path"
	"reflect"
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

	return convertDocument(bytes, out)

}

func convertDocument(data bson.Raw, out runtime.Object) error {

	var resourceVersion int
	if err := data.Lookup("resourceVersion").Unmarshal(&resourceVersion); err != nil {
		logrus.Errorf("unable to get resourceVersion from document: %v", err)
		return err
	}

	raw := data.Lookup("currentRevision")
	if err := raw.Unmarshal(out); err != nil {
		logrus.Errorf("unable to get currentRevision from document: %v", err)
		return err
	}

	var objectID primitive.ObjectID
	if err := objectID.UnmarshalJSON(data.Lookup("_id").Value); err != nil {
		logrus.Errorf("unable to unmarshal objectID from document: %v", err)
		return err
	}

	out.SetUID(objectID.Hex())
	out.SetResourceVersion(resourceVersion)

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

func (s *Store) Watch(ctx context.Context, key string, objPtr runtime.Object, watchFunc func(eventType string, obj runtime.Object)) error {

	keyParts := strings.Split(key, "/")
	group := keyParts[0]
	resource := keyParts[1]

	stream, err := s.mongoClient.Database("core").Collection(group+"/"+resource).Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	go func() {
		for {
			if stream.Next(ctx) {
				operationType := stream.Current.Lookup("operationType")
				operationTypeStr := ""
				if err := operationType.Unmarshal(&operationTypeStr); err != nil {
					logrus.Errorf("unable to unmarshal operationType: %v", err)
					continue
				}

				document := stream.Current.Lookup("fullDocument")
				out := reflect.New(reflect.TypeOf(objPtr).Elem()).Interface().(runtime.Object)
				if err := convertDocument(document.Value, out); err != nil {
					logrus.Errorf("error converting document: %v", err)
					continue
				}

				watchFunc(operationTypeStr, out)
			}
		}
	}()

	<-ctx.Done()

	return nil

}
