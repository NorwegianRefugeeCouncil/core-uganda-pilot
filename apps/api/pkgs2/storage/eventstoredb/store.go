package eventstoredb

import (
	"context"
	"encoding/json"
	"fmt"
	goes "github.com/jetbasrawi/go.geteventstore"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime"
	"github.com/nrc-no/core/apps/api/pkgs2/storage"
	"k8s.io/apimachinery/pkg/conversion"
)

type Store struct {
	client *goes.Client
	new    func() interface{}
	codec  runtime.Codec
}

var _ storage.Interface = &Store{}

func (s Store) Create(ctx context.Context, key string, obj, out runtime.Object) error {
	evt := goes.NewEvent(goes.NewUUID(), "created", obj, nil)
	sw := s.client.NewStreamWriter(key)
	expectedVersion := -1
	if err := sw.Append(&expectedVersion, evt); err != nil {
		return err
	}
	objBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(objBytes, out); err != nil {
		return err
	}
	return nil
}

func (s Store) Get(ctx context.Context, key string, out runtime.Object) error {

	sr := s.client.NewStreamReader(key)
	hasNext := sr.Next()

	if !hasNext {
		return fmt.Errorf("not found")
	}

	var temp = map[string]interface{}{}
	if err := sr.Scan(&temp, nil); err != nil {
		return err
	}

	if sr.Err() != nil {
		return sr.Err()
	}

	evt := sr.EventResponse()
	bytes, err := json.Marshal(temp)
	if err != nil {
		return err
	}

	return decode(s.codec, bytes, out, int64(evt.Event.EventNumber))
}

func (s Store) Update(ctx context.Context, key string, objType runtime.Object, updateFunc storage.UpdateFunc) error {
	panic("implement me")
}

func decode(codec runtime.Codec, value []byte, objPtr runtime.Object, rev int64) error {
	if _, err := conversion.EnforcePtr(objPtr); err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}

	_, _, err := codec.Decode(value, nil, objPtr)
	if err != nil {
		return err
	}

	return nil
}
