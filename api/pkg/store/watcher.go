package store

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"github.com/nrc-no/core/api/pkg/fields"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"reflect"
	"sync"
	"time"
)

const (
	incomingBufSize = 100
	outgoingBufSize = 100
)

type watcher struct {
	database   string
	client     *mongo.Client
	codec      runtime.Codec
	newFunc    func() runtime.Object
	objectType string
}

type watchChan struct {
	watcher           *watcher
	key               string
	initialRev        int64
	recursive         bool
	progressNotify    bool
	ctx               context.Context
	cancelFunc        context.CancelFunc
	incomingEventChan chan *event
	resultChan        chan watch.Event
	errChan           chan error
	selector          fields.Selector
	sentCount         int64
	sentLock          sync.RWMutex
	limit             *int64
}

type event struct {
	key              string
	value            []byte
	prevValue        []byte
	rev              int64
	isDeleted        bool
	isCreated        bool
	isProgressNotify bool
}

func newWatcher(
	client *mongo.Client,
	codec runtime.Codec,
	newFunc func() runtime.Object,
	database string,
) *watcher {
	ret := &watcher{
		client:   client,
		codec:    codec,
		newFunc:  newFunc,
		database: database,
	}
	if newFunc == nil {
		ret.objectType = "<unknown>"
	} else {
		ret.objectType = reflect.TypeOf(newFunc()).String()
	}
	return ret
}

func (w *watcher) Watch(
	ctx context.Context,
	key string,
	rev int64,
	recursive,
	progressNotify bool,
	selector fields.Selector,
	limit *int64) (*watchChan, error) {
	wc := w.createWatchChan(ctx, key, rev, recursive, progressNotify, selector, limit)
	go wc.run()
	return wc, nil
}

func (w *watcher) createWatchChan(ctx context.Context, key string, rev int64, recursive bool, progressNotify bool, selector fields.Selector, limit *int64) *watchChan {
	ctx, cancelFunc := context.WithCancel(ctx)
	wc := &watchChan{
		watcher:           w,
		key:               key,
		initialRev:        rev,
		recursive:         recursive,
		progressNotify:    progressNotify,
		ctx:               ctx,
		cancelFunc:        cancelFunc,
		incomingEventChan: make(chan *event, incomingBufSize),
		resultChan:        make(chan watch.Event, outgoingBufSize),
		errChan:           make(chan error, 1),
		selector:          selector,
		limit:             limit,
	}
	return wc
}

func (wc *watchChan) run() {

	// If the initialRev is 0, then this method will send ADDED
	// events for each item in the store.
	// Then, it will start live-processing of new events
	watchClosedCh := make(chan struct{})
	go wc.startWatching(watchClosedCh)

	// This part receives events emitted by
	// startWatching and pushes them to the result channel
	var resultChanWG sync.WaitGroup
	resultChanWG.Add(1)
	go wc.processEvent(&resultChanWG)

	// Here, we check if we received some error
	// and will basically just exit everything
	select {
	case err := <-wc.errChan:
		if err == context.Canceled {
			break
		}
		errResult := transformErrorToEvent(err)
		if errResult != nil {
			select {
			case <-wc.ctx.Done():
				break
			case wc.resultChan <- *errResult:
				logrus.Errorf("fail")
			}
		}
	case <-wc.ctx.Done():
		break
	case <-watchClosedCh:
	}

	wc.cancelFunc()
	resultChanWG.Wait()
	close(wc.resultChan)
}

func (wc *watchChan) Stop() {
	wc.cancelFunc()
}

func (wc *watchChan) ResultChan() <-chan watch.Event {
	return wc.resultChan
}

func (wc *watchChan) getMongoCollectionWatchFilter(currentOnly bool, fromRev *int64) (bson.M, error) {
	baseFilter := bson.A{}

	if currentOnly {
		baseFilter = append(baseFilter, bson.M{CurrentValueKey: bson.M{"$eq": true}})
	}
	if fromRev != nil {
		fromDate := time.Unix(0, *fromRev)
		baseFilter = append(baseFilter, bson.M{TimestampKey: bson.M{"$gte": fromDate}})
	}

	if !wc.recursive {
		objInfo, err := getMongoObjectInfo(wc.key)
		if err != nil {
			return nil, err
		}
		baseFilter = append(baseFilter, bson.M{Key: bson.M{"$gte": objInfo.key}})
	}

	filter, err := convertFieldSelectorToMongoFilter(wc.selector, fmt.Sprintf("%s.", CurrentValueKey))
	if err != nil {
		return nil, err
	}
	if filter != nil {
		combined := bson.A{
			filter,
		}
		combined = append(combined, baseFilter...)
		return bson.M{

			"$and": combined,
		}, nil
	} else {
		if len(baseFilter) == 0 {
			return nil, nil
		}
		return bson.M{"$and": baseFilter}, nil
	}
}

func (wc *watchChan) getMongoOpLogFilter() (bson.M, error) {

	baseFilter := bson.M{
		fmt.Sprintf("fullDocument.%s", IsCurrentKey): true,
		"operationType": "insert",
	}

	if !wc.recursive {
		objInfo, err := getMongoObjectInfo(wc.key)
		if err != nil {
			return nil, err
		}
		baseFilter[Key] = objInfo.key
	}

	filter, err := convertFieldSelectorToMongoFilter(wc.selector, fmt.Sprintf("fullDocument.%s.", CurrentValueKey))
	if err != nil {
		return nil, err
	}
	if filter != nil {
		filter = bson.M{
			"$and": bson.A{
				baseFilter,
				filter,
			},
		}
	} else {
		filter = baseFilter
	}
	return filter, nil
}

// sync is a method that will retrieve the current state of the store
// and emit an ADDED event for each (non-deleted) object in the store.
// That allows controllers and watchers to get a current state
// snapshot before live-processing the change events
func (wc *watchChan) sync() (int64, error) {

	revision := time.Now().UTC().UnixNano()

	// Recursive means that we're watching multiple objects
	if wc.recursive {

		// Getting the collection name from the provided key
		collInfo, err := getMongoCollectionInfo(wc.key)
		if err != nil {
			return 0, err
		}

		filter, err := wc.getMongoCollectionWatchFilter(true, nil)
		if err != nil {
			return 0, err
		}

		// Find the current items in the store
		findResult, err := wc.watcher.client.
			Database(wc.watcher.database).
			Collection(collInfo.collectionName).
			Find(wc.ctx, filter)
		if err != nil {
			return 0, err
		}

		var foundCount int

		// Process each item
		for {
			if findResult.Next(wc.ctx) {
				cur := findResult.Current

				// Create an ADDED event for each item
				evt, err := parseMongoDataAsEvt(cur)
				if err != nil {
					return 0, err
				}

				foundCount++
				// Send the event to the resultChan
				wc.sendEvent(evt)

			}
			if findResult.Err() != nil {
				return 0, findResult.Err()
			}
			logrus.Infof("found %d records while syncing watch", foundCount)

			break
		}

	} else {

		// This is not a recursive watch, we're only
		// watching a single object

		objInfo, err := getMongoObjectInfo(wc.key)
		if err != nil {
			return 0, err
		}

		filter, err := wc.getMongoCollectionWatchFilter(true, nil)
		if err != nil {
			return 0, err
		}

		// Find the object in the store (the current state)
		findResult := wc.watcher.client.
			Database(wc.watcher.database).
			Collection(objInfo.collection).
			FindOne(wc.ctx, filter)
		if findResult.Err() != nil {
			return 0, findResult.Err()
		}
		raw, err := findResult.DecodeBytes()
		if err != nil {
			return 0, err
		}

		// Send an ADDED event for that entry
		evt, err := parseMongoDataAsEvt(raw)
		if err != nil {
			return 0, err
		}

		// Send the event to resultChan
		wc.sendEvent(evt)
	}
	return revision, nil
}

// parseMongoDataAsEvt is a method that takes in a mongo search result
// and returns an ADDED event corresponding to that object
// This is used by sync to send the initial ADDED events when
// the specified watch initialRev=0
func parseMongoDataAsEvt(res bson.Raw) (*event, error) {

	// We unmarshal the mongo record from the BSON
	record := MongoRecord{}
	if err := bson.Unmarshal(res, &record); err != nil {
		return nil, err
	}

	// We decode the currentValue. We're only interested in the
	// currentValue and not the previousValue. Why? Because the
	// previousValue is only needed for delete/update events.
	// In this case, we know that we're emitting an ADDED event
	var err error
	if record.CurrentValue != nil && len(record.CurrentValue) > 0 {
		record.CurrentBytes, err = json.Marshal(record.CurrentValue)
	}
	if err != nil {
		return nil, err
	}
	if record.PreviousValue != nil && len(record.PreviousValue) > 0 {
		record.PreviousBytes, err = json.Marshal(record.PreviousValue)
	}
	if err != nil {
		return nil, err
	}

	// This is the ADDED event we want to build from the current
	// state of the object
	evt := &event{
		key:       record.Key,
		value:     record.CurrentBytes,
		prevValue: record.PreviousBytes,
		// We're using the record timestamp as a revision
		// This is now extraordinary, but it works.
		// If anyone has a better idea that doesn't include
		// global atomic counters, feel free to contribute.
		rev:              record.Timestamp,
		isDeleted:        record.IsDeleted,
		isCreated:        record.IsCreated,
		isProgressNotify: false,
	}

	return evt, nil
}

// startWatching is a method that
// 1st. If the initialRev=0, is going to emit an ADDED event for every
//      current object in the store. That means that a subscriber would
//      receive an ADDED event for all the current objects. That allows
//      controllers and informers to build a local cache of the current state
//      of the store. If the initialRev!=0, it will not retrieve the current
//      state of the store, and will simply start sending events starting
//      from the given timestamp.
//      The initialRev is actually a timestamp.UnixNano()
// 2nd. It will start watching for the mongo changeStream, will parse
//      these events and convert them to watch.Event, then send them
//      to the result channel
//
// This one is a bit tricky
//
// The reason is that multiple events are going to be emitted for the same update
//
// The current mechanisms for updating an object are illustrated below
//
// ------------------time------------------------>
//
// key ABC     1--------2--------3--------4---------
//             insert   update   delete   added
//
// 1. On new inserts, a single "insert" operation will be emitted from the mongo.ChangeStream
//
// 2. On updates
//    1. A mongo.ChangeStream "update" event will be emitted as we modify the object
//       IsCurrentKey, NextTimestampKey, NextIDKey
//    2. A mongo.ChangeStream "insert" event will be emitted as we insert a new object
//       representing the new version (version 2)
//
// 3. On Deletes
//    1. A mongo.ChangeStream "update" event will be emitted as we modify the object
//       IsDeletedKey
//
// 4. On Inserts after a Delete (same key)
//    1. A mongo.ChangeStream "insert" event will be emitted as we insert a new
//       record representing the new version (version 4)
//    2. A mongo.ChangeStream "update" event will be emitted as we modify version 3
//       IsCurrentKey, NextTimestampKey and NextIDKey
//
// The complexity emerges when we want to start watching at a specific revision. For example,
// a subscriber might maintain some external store and uses Watch to receive new events.
// It also stores until which point (revision) it has successfully reconciled. If the subscriber
// is disconnected or exists, then it can still continue watching events from a certain point in time
// (replay).
//
//
func (wc *watchChan) startWatching(watchClosedCh chan struct{}) {
	defer close(watchClosedCh)

	var hasInitialRev = true

	// If the initialRev = 0, we first get the current state of the
	// store and send ADDED events for each existing object
	if wc.initialRev == 0 {
		hasInitialRev = false
		rev, err := wc.sync()
		if err != nil {
			logrus.Warnf("failed to sync with latest state: %v", err)
			wc.sendError(err)
			return
		}
		wc.initialRev = rev
	}

	var collectionName string

	// If it is not recursive, that means we're only going to be watching
	// a single object
	if !wc.recursive {

		// Get the mongo object info (collection + key) from the provided
		// key
		objInfo, err := getMongoObjectInfo(wc.key)
		if err != nil {
			logrus.Errorf("failed to get collection info from key: %v", err)
			wc.sendError(err)
			return
		}

		collectionName = objInfo.collection
	} else {

		// The watch is recursive, meaning that we're going to watch multiple
		// objects from the store.
		// TODO: allow database filtering based on labels or something
		// Perhaps we don't want to retrieve every single object from the store

		// Get the collection info from the key
		collInfo, err := getMongoCollectionInfo(wc.key)
		if err != nil {
			logrus.Errorf("failed to get collection info from key: %v", err)
			wc.sendError(err)
			return
		}
		collectionName = collInfo.collectionName
	}

	if hasInitialRev {

		logrus.Infof("querying with rv : %d", wc.initialRev)

		filter, err := wc.getMongoCollectionWatchFilter(false, &wc.initialRev)
		if err != nil {
			logrus.Errorf("failed to create collection watch filter: %v", err)
			wc.sendError(err)
			return
		}

		var findOptions []*mongooptions.FindOptions
		if wc.limit != nil {
			wc.sentLock.RLock()
			limit := *wc.limit - wc.sentCount
			wc.sentLock.RUnlock()

			if limit < 0 {
				return
			}
			findOptions = append(findOptions, mongooptions.Find().SetLimit(limit))

		}

		findResult, err := wc.watcher.client.
			Database(wc.watcher.database).
			Collection(collectionName, mongooptions.Collection().SetReadConcern(readconcern.Linearizable())).
			Find(wc.ctx, filter, findOptions...)

		if err != nil {
			logrus.Errorf("watchChannel failed to query collection: %v", err)
			wc.sendError(err)
			return
		}

		var foundQry int64 = 0
		for {
			if findResult.Next(wc.ctx) {
				cur := findResult.Current
				var mongoRecord MongoRecord
				if err := bson.Unmarshal(cur, &mongoRecord); err != nil {
					logWatchChannelErr(err)
					wc.sendError(err)
					return
				}
				ev, err := convertMongoDocumentToEvent(mongoRecord)
				if err != nil {
					logWatchChannelErr(err)
					wc.sendError(err)
					return
				}
				foundQry++
				wc.sendEvent(ev)
				wc.initialRev = ev.rev + 1
			}
			if findResult.Err() != nil {
				logWatchChannelErr(err)
				wc.sendError(err)
				return
			}
			break
		}

		logrus.Infof("found %d events using query", foundQry)

	}

	// If the provided initialRev is not 0, then start watching
	// events at the given timestamp
	var watchOptions = []*mongooptions.ChangeStreamOptions{
		mongooptions.ChangeStream().SetFullDocument(mongooptions.UpdateLookup),
	}
	if wc.initialRev != 0 {
		ts := time.Unix(0, wc.initialRev).UTC().Unix()
		watchOptions = append(watchOptions, mongooptions.ChangeStream().SetStartAtOperationTime(&primitive.Timestamp{
			T: uint32(ts),
		}))
	}

	filter, err := wc.getMongoOpLogFilter()
	if err != nil {
		logrus.Warnf("failed to create mongo filter: %v", err)
		wc.sendError(err)
		return
	}
	var pipeline interface{}
	if filter != nil {
		pipeline = mongo.Pipeline{bson.D{{"$match", filter}}}
	} else {
		pipeline = mongo.Pipeline{}
	}

	// Start watching the mongo ChangeStream
	cs, err := wc.watcher.client.
		Database(wc.watcher.database).
		Collection(collectionName).
		Watch(wc.ctx, pipeline, watchOptions...)
	if err != nil {
		logWatchChannelErr(err)
		wc.sendError(err)
		return
	}

	// This is a blocking loop. cs.Next will block until the
	// next result is available
	// TODO: add progressNotify every n seconds by sending
	// a BOOKMARK event meaning that the connection is not dead,
	// but just not receiving new events. Make sure to test
	// that the connection is alive using some kind of PING
	// for the BOOKMARK event
	for {
		if cs.Next(wc.ctx) {

			// We convert the mongo ChangeStream event
			// to the *event type
			evt, err := parseEvent(cs.Current)
			if err != nil {
				logWatchChannelErr(err)
				wc.sendError(err)
				return
			}

			if evt.rev < wc.initialRev {
				continue
			}

			// The returned event might be nil
			if evt == nil {
				continue
			}

			// Send the event to the incoming event channel
			// for further processing
			wc.sendEvent(evt)
		}

		// At this point, either the context was canceled
		// or we have a mongo error. Perhaps the connection
		// was closed
		if cs.Err() != nil {
			if errors2.Is(cs.Err(), context.Canceled) {
				return
			}
			logWatchChannelErr(cs.Err())
			wc.sendError(cs.Err())
			return
		}
	}
}

// processEvents receives events from the incomingEventChan and
// converts them to a watch.Event, then sends them to the
// resultChan
func (wc *watchChan) processEvent(resultChanWG *sync.WaitGroup) {
	defer resultChanWG.Done()
	for {
		select {
		case e := <-wc.incomingEventChan:
			res := wc.transform(e)
			if res == nil {
				continue
			}
			if len(wc.resultChan) == outgoingBufSize {
				logrus.Warnf("fast watcher, slow processing")
			}

			select {
			case wc.resultChan <- *res:
				wc.sentLock.Lock()
				wc.sentCount++
				wc.sentLock.Unlock()

				if wc.limit != nil && wc.sentCount == *wc.limit {
					return
				}

			case <-wc.ctx.Done():
				return
			}
		case <-wc.ctx.Done():
			return
		}
	}
}

// sendError sends an error to the error channel if the context is not already done
func (wc *watchChan) sendError(err error) {
	select {
	case wc.errChan <- err:
	case <-wc.ctx.Done():
		return
	}
}

// sendEvent sends an event to the incomingEventChannel, so that it can get converted
// to a watch.Event and then sent to the result channel
func (wc *watchChan) sendEvent(evt *event) {
	if len(wc.incomingEventChan) == incomingBufSize {
		logrus.Warnf("fast watcher, slow processing")
	}
	select {
	case wc.incomingEventChan <- evt:
	case <-wc.ctx.Done():
		return
	}
}

// transform transforms an event into a watch.Event to send to the result channel
func (wc *watchChan) transform(e *event) (res *watch.Event) {
	curObj, oldObj, err := wc.prepareObjs(e)
	if err != nil {
		logrus.Errorf("failed to prepare current and previous objects: %v", err)
		wc.sendError(err)
		return nil
	}
	switch {
	case e.isProgressNotify:
		obj := wc.watcher.newFunc()
		// todo: update object version
		res = &watch.Event{
			Type:   watch.Bookmark,
			Object: obj,
		}
	case e.isDeleted:
		res = &watch.Event{
			Type:   watch.Deleted,
			Object: oldObj,
		}
	case e.isCreated:
		res = &watch.Event{
			Type:   watch.Added,
			Object: curObj,
		}
	default:
		// TODO: emit ADDED if the modified object causes it to actually pass the filter but the previous one did not
		res = &watch.Event{
			Type:   watch.Modified,
			Object: curObj,
		}
	}
	return res
}

type MongoUpdateDescription struct {
	UpdatedFields map[string]interface{} `bson:"updatedFields"`
	RemovedFields []string               `bson:"removedFields"`
}

func (m MongoUpdateDescription) IsDeleted() bool {
	if m.UpdatedFields == nil {
		return false
	}
	val, ok := m.UpdatedFields[IsDeletedKey]
	if !ok {
		return false
	}
	valBool, ok := val.(bool)
	if !ok {
		return false
	}
	return valBool
}

type MongoChangeStreamEvt struct {
	FullDocument      MongoRecord            `bson:"fullDocument"`
	UpdateDescription MongoUpdateDescription `bson:"updateDescription"`
	OperationType     string                 `bson:"operationType"`
}

func parseEvent(current bson.Raw) (*event, error) {
	var changeStreamEvt = MongoChangeStreamEvt{}
	if err := bson.Unmarshal(current, &changeStreamEvt); err != nil {
		return nil, err
	}

	evt, err := convertMongoDocumentToEvent(changeStreamEvt.FullDocument)
	if err != nil {
		return nil, err
	}

	return evt, nil
}

func convertMongoDocumentToEvent(document MongoRecord) (*event, error) {

	var value []byte
	var prevValue []byte
	var err error

	if document.CurrentValue != nil && len(document.CurrentValue) > 0 {
		document.CurrentBytes, err = json.Marshal(document.CurrentValue)
	}
	if err != nil {
		return nil, err
	}

	if document.PreviousValue != nil && len(document.PreviousValue) > 0 {
		document.PreviousBytes, err = json.Marshal(document.PreviousValue)
	}
	if err != nil {
		return nil, err
	}

	if document.IsDeleted {
		prevValue = document.PreviousBytes
	} else if document.IsCreated {
		value = document.CurrentBytes
	} else {
		prevValue = document.PreviousBytes
		value = document.CurrentBytes
	}

	return &event{
		key:              document.Key,
		value:            value,
		prevValue:        prevValue,
		rev:              document.Timestamp,
		isDeleted:        document.IsDeleted,
		isCreated:        document.IsCreated,
		isProgressNotify: false,
	}, nil

}

// prepareObjs retrieves objects from the event payload
func (wc *watchChan) prepareObjs(e *event) (curObj, oldObj runtime.Object, err error) {

	if e.isProgressNotify {
		return nil, nil, nil
	}

	if !e.isDeleted {
		curObj, err = decodeEventObj(wc.watcher.codec, e.value, e.rev)
		if err != nil {
			logrus.Errorf("unable to decode current value: %v", err)
			return nil, nil, err
		}
	}

	if len(e.prevValue) > 0 && (e.isDeleted) {
		oldObj, err = decodeEventObj(wc.watcher.codec, e.prevValue, e.rev)
		if err != nil {
			logrus.Errorf("unable to decode previous value: %v", err)
			return nil, nil, err
		}
	}
	return curObj, oldObj, nil

}

func decodeEventObj(codec runtime.Codec, data []byte, rev int64) (runtime.Object, error) {
	obj, err := runtime.Decode(codec, data)
	if err != nil {
		logrus.Errorf("unable to decode event object: %v", err)
		return nil, err
	}
	return obj, nil
}

func logWatchChannelErr(err error) {
	logrus.Warnf("watch chan error: %v", err)
}

func transformErrorToEvent(err error) *watch.Event {
	err = interpretWatchError(err)
	if _, ok := err.(errors.APIStatus); !ok {
		err = errors.NewInternalError(err)
	}
	status := err.(errors.APIStatus).Status()
	return &watch.Event{
		Type:   watch.Error,
		Object: &status,
	}
}

func interpretWatchError(err error) error {
	// TODO: interpret mongo changestream errors
	return err
}
