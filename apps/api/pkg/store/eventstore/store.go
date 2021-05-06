package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	goes "github.com/jetbasrawi/go.geteventstore"
	"github.com/nrc-no/core/apps/api/pkg/store"
	"github.com/nrc-no/core/apps/api/pkg/util/pointers"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"path"
	"strconv"
	"time"
)

type EventType string

const (
	Created           EventType = "Created"
	Updated           EventType = "Updated"
	MarkedForDeletion EventType = "MarkedForDeletion"
	Deleted           EventType = "Deleted"
)

type Store struct {
	client *goes.Client
	codec  runtime.Codec
}

var _ store.Interface = &Store{}

func NewStore(goesClient *goes.Client, codec runtime.Codec) *Store {
	return &Store{
		client: goesClient,
		codec:  codec,
	}
}

type StoreEvent struct {
}

func getEventTypePrefix(gvk schema.GroupVersionKind) string {

	// nrc.core.no::MyCustomForm::Created
	// nrc.core.no::MyCustomForm::Updated
	// nrc.core.no::MyCustomForm::MarkedForDeletion
	// nrc.core.no::MyCustomForm::Deleted

	group, kind := gvk.Group, gvk.Kind

	if len(group) == 0 {
		return fmt.Sprintf("%s", kind)
	} else {
		return fmt.Sprintf("%s::%s", group, kind)
	}
}

func getEventType(gvk schema.GroupVersionKind, eventType EventType) string {
	return fmt.Sprintf("%s::%s", getEventTypePrefix(gvk), eventType)
}

func (s Store) Create(ctx context.Context, key string, obj, out runtime.Object) error {

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetCreationTimestamp(v1.NewTime(time.Now().UTC()))

	data, err := runtime.Encode(s.codec, obj)
	if err != nil {
		return err
	}

	gvk := obj.GetObjectKind().GroupVersionKind()
	key = path.Join(key)

	// We need a temporary structure, because eventstore client does not accept
	// raw bytes as a body. It only accepts an object that it will then serialize.
	// Also, we need to rely on k8s-apimachinery codec. Hence the need of a temporary
	// structure.
	var tempStruct = map[string]interface{}{}
	if err := json.Unmarshal(data, &tempStruct); err != nil {
		return fmt.Errorf("unable to convert data to json: %v", err)
	}

	evt := goes.NewEvent(goes.NewUUID(), getEventType(gvk, Created), tempStruct, nil)
	sw := s.client.NewStreamWriter(key)

	// An expected value of -1 means that the stream should not exist at the time of writing.
	// The stream will be created
	if err := sw.Append(pointers.NewIntPtr(-2), evt); err != nil {
		return fmt.Errorf("unable to append event to store: %v", err)
	}

	if out != nil {
		return decode(s.codec, data, out, 1)
	}

	return nil
}

func (s Store) Get(ctx context.Context, key string, getOptions store.GetOptions, objPtr runtime.Object) error {

	sr := s.client.NewStreamReader(key)
	tempStr := map[string]interface{}{}
	hasNext := sr.Next()

	if !hasNext {
		if getOptions.IgnoreNotFound {
			return runtime.SetZeroValue(objPtr)
		}
		return fmt.Errorf("not found")
	}

	er := sr.EventResponse()
	if err := sr.Scan(&tempStr, nil); err != nil {
		return fmt.Errorf("could not decode eventstore event: %v", err)
	}

	bytes, err := json.Marshal(&tempStr)
	if err != nil {
		return fmt.Errorf("could not convert eventstore event to bytearray: %v", err)
	}

	return decode(s.codec, bytes, objPtr, uint64(er.Event.EventNumber))

}

func (s Store) Update(ctx context.Context, key string, objType runtime.Object, updateFunc store.UpdateFunc, updateOptions store.UpdateOptions) error {
	panic("implement me")
}

func decode(codec runtime.Codec, value []byte, objPtr runtime.Object, rev uint64) error {
	if _, err := conversion.EnforcePtr(objPtr); err != nil {
		return fmt.Errorf("unable to convert output object to pointer: %v", err)
	}
	_, _, err := codec.Decode(value, nil, objPtr)
	if err != nil {
		return err
	}
	accessor, err := meta.Accessor(objPtr)
	if err != nil {
		return err
	}
	versionString := ""
	if rev != 0 {
		versionString = strconv.FormatUint(rev, 10)
	}
	accessor.SetResourceVersion(versionString)
	return nil
}
