package store

import (
	"context"
	"fmt"
	"github.com/nrc-no/core/api/pkg/apis/core"
	corev1 "github.com/nrc-no/core/api/pkg/apis/core/v1"
	"github.com/nrc-no/core/api/pkg/fields"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	ctx          context.Context
	cancel       context.CancelFunc
	client       *mongo.Client
	scheme       *runtime.Scheme
	codecFactory serializer.CodecFactory
	codec        runtime.Codec
	store        *MongoStore
	databaseName string
}

func (s *Suite) TestCRUD() {

	name := uuid.NewV4().String()

	fd := &corev1.FormDefinition{
		TypeMeta: metav1.TypeMeta{
			Kind:       "FormDefinition",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	var evtChan = make(chan struct{})

	go func() {
		w, err := s.store.Watch(s.ctx, getListKey(), ListOptions{
			ResourceVersion: strconv.FormatInt(time.Now().UTC().UnixNano(), 10),
		})
		if err != nil {
			s.T().Fatal(err)
		}
		var i int = 0
		for e := range w.ResultChan() {
			i++
			s.T().Logf("%s - %#v", e.Type, e.Object)
			if i == 4 {
				evtChan <- struct{}{}
			}
		}
	}()

	// Test creating the form definition
	// We should be able to store new items in the store
	out := &corev1.FormDefinition{}
	key := getObjectKey(fd)
	if err := s.store.Create(s.ctx, key, CreateOptions{}, fd, out); !assert.NoError(s.T(), err) {
		return
	}

	// Test listing the form definitions
	// That would return a list of all the FormDefinitions
	list := &corev1.FormDefinitionList{}
	if err := s.store.List(s.ctx, getListKey(), ListOptions{}, list); !assert.NoError(s.T(), err) {
		return
	}
	var found = false
	for _, item := range list.Items {
		if item.Name == fd.Name {
			found = true
			break
		}
	}
	assert.True(s.T(), found)

	// Test listing with a predicate
	// We filter the predicate by the metadata.name field,
	// since it's a uuid generated on every test run
	// it should return only a single result
	filteredList := &corev1.FormDefinitionList{}
	if err := s.store.List(s.ctx, getListKey(), ListOptions{
		Selector: fields.Equal("metadata.name", name),
	}, filteredList); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), 1, len(filteredList.Items))

	// Test listing with a predicate that should not match
	// We still use the uuid + "a", so that we know that there
	// won't be any record matching this
	filteredList2 := &corev1.FormDefinitionList{}
	if err := s.store.List(s.ctx, getListKey(), ListOptions{
		Selector: fields.Equal("metadata.name", name+"a"),
	}, filteredList2); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), 0, len(filteredList2.Items))

	// Test getting back the form definition
	// We should be able to execute simple "gets"
	in := &corev1.FormDefinition{}
	if err := s.store.Get(s.ctx, key, GetOptions{}, in); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), fd, in)

	// Test updating the form definition
	// We should be able to update the form definitions
	updated := &corev1.FormDefinition{}
	if err := s.store.Update(s.ctx, key, UpdateOptions{}, updated, func(input runtime.Object) (output runtime.Object, err error) {
		formDef := input.(*corev1.FormDefinition)
		formDef.Spec.Names.Kind = "Updated"
		return formDef, nil
	}); !assert.NoError(s.T(), err) {
		return
	}
	expected := fd.DeepCopy()
	expected.Spec.Names.Kind = "Updated"
	assert.Equal(s.T(), expected, updated)

	// Test deleting the form definition
	// We should also be able to delete the form definitions
	deleted := &corev1.FormDefinition{}
	if err := s.store.Delete(s.ctx, key, DeleteOptions{}, deleted); !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), expected, deleted)

	// Test getting the form definition again
	// That should return an error with "no documents in result"
	afterDelete := &corev1.FormDefinition{}
	err := s.store.Delete(s.ctx, key, DeleteOptions{}, afterDelete)
	if !assert.Error(s.T(), err) {
		return
	}
	assert.Equal(s.T(), "mongo: no documents in result", err.Error())

	// Test creating a form definition with the same name again
	recreated := &corev1.FormDefinition{}
	if err := s.store.Create(s.ctx, key, CreateOptions{}, fd, recreated); !assert.NoError(s.T(), err) {
		return
	}

	// Test recreating a form definition with the same name again
	// should throw a key conflict
	recreatedConflict := &corev1.FormDefinition{}
	if err := s.store.Create(s.ctx, key, CreateOptions{}, fd, recreatedConflict); !assert.Error(s.T(), err) {
		return
	}

	<-evtChan

}

func TestConvertSelectorToMongoFilter(t *testing.T) {

	var selector = fields.And(
		fields.Or(
			fields.Equal("Key1", "Value1"),
			fields.Equal("Key2", 123),
		),
		fields.Equal("Key3", 10.2),
	)

	expr, err := convertFieldSelectorToMongoFilter(selector, CurrentValueKey+".")
	if !assert.NoError(t, err) {
		return
	}

	expected := bson.M{
		"$and": bson.A{
			bson.M{
				"$or": bson.A{
					bson.M{"currentValue.Key1": bson.M{"$eq": "Value1"}},
					bson.M{"currentValue.Key2": bson.M{"$eq": 123}},
				},
			},
			bson.M{"currentValue.Key3": bson.M{"$eq": 10.2}},
		},
	}

	assert.Equal(t, expected, expr)

}

func (s *Suite) TestWatch() {

	// Timeout for the test is 5 seconds
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*1)
	defer cancel()

	// A basic form definition
	formDef := &corev1.FormDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testwatch",
		},
	}
	key := getObjectKey(formDef)

	// Cleanup in case an object with the same name already exists
	if !ensureDeleted(s.T(), s.ctx, s.store, formDef) {
		return
	}

	// Start recording events for that specific form definition
	r, err := newEventRecorder(s.T(), ctx, s.store, key, nameListOptions(formDef.Name))
	if !assert.NoError(s.T(), err) {
		return
	}

	// Create the form definition
	created := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	// Assert that an ADDED event was issued
	r.AssertNextEventOrDie(s.T(), isAddedEvent())

	// Update the form definition
	if !assert.NoError(s.T(), s.store.Update(s.ctx, key, UpdateOptions{}, formDef, func(input runtime.Object) (output runtime.Object, err error) {
		input.(*corev1.FormDefinition).Spec.Group = "Updated"
		return input, nil
	})) {
		return
	}

	r.AssertNextEventOrDie(s.T(), isModifiedEvent())

	// Delete the form definition
	deleted := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Delete(s.ctx, key, DeleteOptions{}, &deleted)) {
		return
	}

	r.AssertNextEventOrDie(s.T(), isDeletedEvent())

	// Create the form definition again on the same key
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	r.AssertNextEventOrDie(s.T(), isAddedEvent())

}

func (s *Suite) TestWatchExistingResource() {
	// Timeout for the test is 5 seconds
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*1)
	defer cancel()

	// A basic form definition
	formDef := &corev1.FormDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testwatch",
		},
	}
	key := getObjectKey(formDef)

	// Cleanup in case an object with the same name already exists
	if !ensureDeleted(s.T(), s.ctx, s.store, formDef) {
		return
	}

	// Create the form definition
	created := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	// Update the form definition
	if !assert.NoError(s.T(), s.store.Update(s.ctx, key, UpdateOptions{}, formDef, func(input runtime.Object) (output runtime.Object, err error) {
		input.(*corev1.FormDefinition).Spec.Group = "Updated"
		return input, nil
	})) {
		return
	}

	// Delete the form definition
	deleted := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Delete(s.ctx, key, DeleteOptions{}, &deleted)) {
		return
	}

	// Create the form definition again on the same key
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	// Start recording events for that specific form definition
	r, err := newEventRecorder(s.T(), ctx, s.store, key, nameListOptions(formDef.Name))
	if !assert.NoError(s.T(), err) {
		return
	}
	r.AssertNextEventOrDie(s.T(), isAddedEvent())
}

func (s *Suite) TestWatchExistingResourceAtGivenVersion() {
	// Timeout for the test is 5 seconds
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*3)
	defer cancel()

	// A basic form definition
	formDef := &corev1.FormDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "watchgivenversion",
		},
	}
	key := getObjectKey(formDef)

	// Cleanup in case an object with the same name already exists
	if !ensureDeleted(s.T(), s.ctx, s.store, formDef) {
		return
	}

	// Create the form definition
	created := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	// Update the form definition
	var updated corev1.FormDefinition
	if !assert.NoError(s.T(), s.store.Update(s.ctx, key, UpdateOptions{}, &updated, func(input runtime.Object) (output runtime.Object, err error) {
		input.(*corev1.FormDefinition).Spec.Group = "Updated"
		return input, nil
	})) {
		return
	}

	// Delete the form definition
	deleted := corev1.FormDefinition{}
	if !assert.NoError(s.T(), s.store.Delete(s.ctx, key, DeleteOptions{}, &deleted)) {
		return
	}

	// Create the form definition again on the same key
	if !assert.NoError(s.T(), s.store.Create(s.ctx, key, CreateOptions{}, formDef, &created)) {
		return
	}

	logrus.Infof("watching with rv : %s", deleted.ResourceVersion)

	// Start recording events for that specific form definition
	var limit int64 = 2
	r, err := newEventRecorder(s.T(), ctx, s.store, key, ListOptions{
		Selector:        fields.Equal("metadata.name", formDef.Name),
		ResourceVersion: deleted.ResourceVersion,
		Limit:           &limit,
	})
	if !assert.NoError(s.T(), err) {
		return
	}
	r.AssertNextEventOrDie(s.T(), isDeletedEvent())
	r.AssertNextEventOrDie(s.T(), isAddedEvent())
}

func nameListOptions(name string) ListOptions {
	return ListOptions{Selector: fields.Equal("metadata.name", name)}
}

type eventRecorder struct {
	events []watch.Event
	lock   sync.RWMutex
	ctx    context.Context
	idx    int
}

type matchEventFunc func(t *testing.T, evt *watch.Event) bool

func isDeletedEvent() matchEventFunc {
	return eventOfType(watch.Deleted)
}
func isAddedEvent() matchEventFunc {
	return eventOfType(watch.Added)
}
func isModifiedEvent() matchEventFunc {
	return eventOfType(watch.Modified)
}
func isBookmarkEvent() matchEventFunc {
	return eventOfType(watch.Bookmark)
}

func eventOfType(evtType watch.EventType) matchEventFunc {
	return func(t *testing.T, evt *watch.Event) bool {
		if !assert.NotNil(t, evt) {
			return false
		}
		if !assert.Equal(t, evtType, evt.Type) {
			return false
		}
		return true
	}
}

func (r *eventRecorder) NextOrDie(t *testing.T) (evt *watch.Event) {
	return r.AssertNextEventOrDie(t, func(t *testing.T, evt *watch.Event) bool {
		return true
	})
}

func (r *eventRecorder) AssertNextEventOrDie(t *testing.T, pred func(t *testing.T, evt *watch.Event) bool) (evt *watch.Event) {
	r.idx++
	if err := wait.PollImmediateUntilWithContext(r.ctx, time.Millisecond*20, func(ctx context.Context) (done bool, err error) {
		r.lock.RLock()
		defer r.lock.RUnlock()
		if len(r.events) > r.idx {
			evt = &r.events[r.idx]
			if !pred(t, evt) {
				return true, fmt.Errorf("predicate did not match event")
			}
			return true, nil
		}
		return false, nil
	}); err != nil {
		t.Fatalf("could not get next event: %v", err)
	}
	return evt
}

func (r *eventRecorder) GetEvents() []watch.Event {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.events
}

func (r *eventRecorder) WaitForEventCount(count int) error {
	if err := wait.PollImmediateUntilWithContext(r.ctx, time.Millisecond*20, func(ctx context.Context) (done bool, err error) {
		r.lock.RLock()
		defer r.lock.RUnlock()
		if len(r.events) >= count {
			return true, nil
		}
		return false, nil
	}); err != nil {
		return err
	}
	return nil
}

func newEventRecorder(t *testing.T, ctx context.Context, store *MongoStore, key string, listOptions ListOptions) (*eventRecorder, error) {
	r := eventRecorder{
		lock: sync.RWMutex{},
		ctx:  ctx,
		idx:  -1,
	}
	w, err := store.Watch(ctx, key, listOptions)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case evt := <-w.ResultChan():
				r.lock.Lock()
				t.Logf("received new event: %s", evt.Type)
				r.events = append(r.events, evt)
				r.lock.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()
	return &r, nil
}

func ensureDeleted(t *testing.T, ctx context.Context, store *MongoStore, formDef *corev1.FormDefinition) bool {
	key := getObjectKey(formDef)
	deleted := corev1.FormDefinition{}
	if err := store.Delete(ctx, key, DeleteOptions{}, &deleted); err != nil && !IsNotFound(err) {
		assert.NoError(t, err)
		return false
	}
	return true
}

func getObjectKey(fd *corev1.FormDefinition) string {
	return strings.Join([]string{
		corev1.SchemeGroupVersion.Group,
		"formdefinitions",
		fd.Name,
	}, "/")
}

func getListKey() string {
	return strings.Join([]string{
		corev1.SchemeGroupVersion.Group,
		"formdefinitions",
	}, "/")
}

func (s *Suite) SetupSuite() {

	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017,localhost:27018,localhost:27019"),
	//SetAuth(options.Credential{
	//	Username: "root",
	//	Password: "pass12345",
	//}),
	)
	if err != nil {
		s.T().Fatal(err)
	}
	s.client = client

	scheme := runtime.NewScheme()
	s.scheme = scheme
	codecFactory := serializer.NewCodecFactory(scheme)
	s.codecFactory = codecFactory

	if err := corev1.AddToScheme(scheme); err != nil {
		s.T().Fatal(err)
	}
	if err := core.AddToScheme(scheme); err != nil {
		s.T().Fatal(err)
	}

	codec := codecFactory.LegacyCodec(corev1.SchemeGroupVersion)
	s.codec = codec

	s.databaseName = "test"
	store := NewMongoStore(client, codec, s.databaseName, func() runtime.Object {
		return &corev1.FormDefinition{}
	})
	s.store = store

}

func (s *Suite) TearDownSuite() {
	s.cancel()
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
