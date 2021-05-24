package store

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
	"time"
)

type DestroyFunc func()

func Create(c Config, newFunc func() runtime.Object) (Interface, DestroyFunc, error) {

	uri := strings.Join(c.Transport.ServerList, ",")

	storeOptions := []*options.ClientOptions{}
	storeOptions = append(storeOptions, options.Client().ApplyURI(uri))

	if len(c.Transport.Username) > 0 && len(c.Transport.Password) > 0 {
		storeOptions = append(storeOptions, options.Client().SetAuth(
			options.Credential{
				Username: c.Transport.Username,
				Password: c.Transport.Password,
			},
		))
	}

	client, err := mongo.NewClient(
		storeOptions...,
	)
	if err != nil {
		return nil, nil, err
	}
	if err := client.Connect(context.TODO()); err != nil {
		return nil, nil, err
	}

	store := NewMongoStore(client, c.Codec, c.Transport.Database, newFunc)

	return store, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		client.Disconnect(ctx)
	}, nil

}
