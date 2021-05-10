package backend

import (
	"context"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	mongostorage "github.com/nrc-no/core/apps/api/pkg/storage/mongo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

func Create(c Config, newFunc func() runtime.Object) (storage.Interface, func(), error) {
	return NewMongoStore(c, newFunc)
}

func NewMongoStore(c Config, newFunc func() runtime.Object) (storage.Interface, func(), error) {
	clientOptions := mongooptions.Client().ApplyURI("mongodb://" + strings.Join(c.Transport.ServerList, ","))
	mongoClient, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)

	logrus.Info("about to connect to mongodb")
	if err := mongoClient.Connect(ctx); err != nil {
		cancel()
		return nil, nil, err
	}
	logrus.Info("connected to mongodb")

	store, err := mongostorage.NewStore(mongoClient, c.Codec, newFunc, c.Prefix)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	return store, func() {
		logrus.Info("shutting down mongodb client")
		cancel()
	}, nil

}
