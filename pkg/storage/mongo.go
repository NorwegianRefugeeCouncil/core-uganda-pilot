package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hash"
	"os"
	"sync"
)

type MongoClientSrc interface {
	GetMongoClient() (*mongo.Client, error)
}

type fileMongoCLientSrc struct {
	ctx context.Context
	usernameFile string
	passwordFile string
	username string
	password string
	mongoHosts []string
	mongoClient *mongo.Client
	credentialsHash hash.Hash
	lock sync.RWMutex
	fileReader func(fn string) (string, error)
}

func (f *fileMongoCLientSrc) GetMongoClient() (*mongo.Client, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	var err error

	var username []byte
	var password []byte

	if len(f.usernameFile) > 0 {
		if err = readFileInto(f.usernameFile, username); err != nil {
			return nil, err
		}
	} else if len(f.username) > 0 {
		username = []byte(f.username)
	} else {
		return nil, errors.New("missing mongo username")
	}

	if len(f.passwordFile) > 0 {
		if err = readFileInto(f.passwordFile, password); err != nil {
			return nil, err
		}
	} else if len(f.password) > 0 {
		password = []byte(f.password)
	} else {
		return nil, errors.New("missing mongo password")
	}

	var h = sha256.New()

	h.Write(bytes.Join([][]byte{username, password}, nil))

	if f.mongoClient == nil || !bytes.Equal(h.Sum(nil), f.credentialsHash.Sum(nil)) {
		 client, err := MongoClient(f.mongoHosts, string(username), string(password))
		if err != nil {
			return nil, err
		}
		err = client.Connect(f.ctx)
		if err != nil {
			return nil, err
		}
		err = client.Ping(f.ctx, nil)
		if err != nil {
			panic(err)
		}
		f.mongoClient = client
		f.credentialsHash = h
	}

	return f.mongoClient, nil
}

func NewMongoClientSrc(ctx context.Context, usernameFile string, passwordFile string, username string, password string, mongoHosts []string) MongoClientSrc {
	return &fileMongoCLientSrc{
		ctx: ctx,
		usernameFile: usernameFile,
		passwordFile: passwordFile,
		username:     username,
		password:     password,
		mongoHosts:   mongoHosts,
		credentialsHash: sha256.New(),
	}
}

func MongoClient(hosts []string, username, password string) (*mongo.Client, error) {
	// Setup mongo client
	mongoClient, err := mongo.NewClient(options.Client().
		SetHosts(hosts).
		SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			}))
	if err != nil {
		return nil, err
	}
	return mongoClient, nil
}

func readFileInto(filename string, b []byte) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	b = data
	return nil
}

func GetCollectionFn(mongoClientSrc MongoClientSrc, database, collection string) func() (*mongo.Collection, error) {
	return func() (*mongo.Collection, error) {
		client, err := mongoClientSrc.GetMongoClient()
		if err != nil {
			return nil, err
		}
		collection := client.Database(database).Collection(collection)

		return collection, nil
	}
}
