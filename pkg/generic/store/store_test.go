package store

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type dummyObj struct {
	ID      string `bson:"id"`
	Version int    `bson:"version"`
	Name    string `bson:"name"`
}

func (d *dummyObj) SetVersion(version int) {
	d.Version = version
}
func (d *dummyObj) GetVersion() int {
	return d.Version
}

func clientFn() (*mongo.Client, error) {
	return mongo.Connect(
		context.Background(),
		options.Client().SetHosts([]string{
			"localhost:27017",
			"localhost:27018",
			"localhost:27019",
		}))
}

func store() *genericStore {
	return &genericStore{
		databaseName:   "test_genericStore",
		collectionName: "dummyObj",
		clientFn:       clientFn,
	}
}

func Test_genericStore_Create(t *testing.T) {

	type args struct {
		ctx           context.Context
		obj           *dummyObj
		createOptions CreateOptions
	}
	tests := []struct {
		name    string
		args    args
		want    dummyObj
		wantErr bool
	}{
		{
			name: "shouldSetVersion1OnCreate",
			args: args{
				ctx:           context.Background(),
				obj:           &dummyObj{ID: "abc", Version: 0, Name: "bob"},
				createOptions: CreateOptions{},
			},
			want:    dummyObj{ID: "abc", Version: 1, Name: "bob"},
			wantErr: false,
		}, {
			name: "shouldIgnoreVersionOnCreate",
			args: args{
				ctx:           context.Background(),
				obj:           &dummyObj{ID: "create-ignore-version", Version: 0, Name: "bob"},
				createOptions: CreateOptions{},
			},
			want:    dummyObj{ID: "abc", Version: 1, Name: "bob"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store()
			out := &dummyObj{}
			if err := s.Create(tt.args.ctx, tt.args.obj, out, tt.args.createOptions); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, &tt.want, tt.args.obj)

			_ = cleanupObj(tt.args.ctx, s, tt.args.obj)
		})
	}
}

func Test_genericStore_Update(t *testing.T) {

	s := store()
	ctx := context.Background()
	key := "test-key"

	tests := []struct {
		name            string
		key             string
		expectNoUpdate  bool
		wantNotFoundErr bool
		ignoreNotFound  bool
	}{
		{
			// Updating a document should save a new version
			name:            "updateExisting",
			key:             key,
			ignoreNotFound:  false,
			expectNoUpdate:  false,
			wantNotFoundErr: false,
		}, {
			// Updating a non-existing document with ignoreNotFound=false document should throw
			name:            "updateNonExisting",
			key:             "/non-existing",
			ignoreNotFound:  false,
			expectNoUpdate:  false,
			wantNotFoundErr: true,
		}, {
			// Updating a document with the same data should not create a new version
			name:            "updateWithSameData",
			key:             key,
			ignoreNotFound:  false,
			expectNoUpdate:  true,
			wantNotFoundErr: false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			obj := &dummyObj{ID: key, Version: 0, Name: "bob"}
			if err := populateObj(ctx, s, obj); !assert.NoError(t, err) {
				return
			}
			defer cleanupObj(ctx, s, obj)

			oldVersion := obj.GetVersion()

			name := fmt.Sprintf("bob-%d", i)
			if tt.expectNoUpdate {
				name = obj.Name
			}

			out := &dummyObj{}

			err := s.Update(ctx, tt.key, UpdateOptions{IgnoreNotFound: tt.ignoreNotFound}, out, func(in Object) (Object, error) {
				o := in.(*dummyObj)
				if !tt.expectNoUpdate {
					o.Name = fmt.Sprintf("bob-%d", i)
				}
				return o, nil
			})

			if tt.wantNotFoundErr {
				assert.Error(t, err)
				return
			}

			if tt.expectNoUpdate {
				assert.Equal(t, oldVersion, out.GetVersion())
			} else {
				assert.Equal(t, oldVersion+1, out.GetVersion())
			}

			assert.Equal(t, name, out.Name)

		})
	}
}

func populateObj(ctx context.Context, s *genericStore, obj *dummyObj) error {
	if err := cleanupObj(ctx, s, obj); err != nil {
		return err
	}
	if err := s.Create(ctx, obj, nil, CreateOptions{}); err != nil {
		return err
	}
	return nil
}

func cleanupObj(ctx context.Context, s *genericStore, obj *dummyObj) error {
	mongoCli, err := s.clientFn()
	if err != nil {
		return err
	}
	_, err = mongoCli.Database(s.databaseName).Collection(s.collectionName).DeleteMany(ctx, bson.M{
		"id": obj.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
