package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/storage"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
)

const (
	// We have set a buffer in order to reduce times of context switches.
	incomingBufSize = 100
	outgoingBufSize = 100
)

type watcher struct {
	collection *mongo.Collection
	codec      runtime.Codec
	newFunc    func() runtime.Object
	objectType string
	versioner  storage.Versioner
}

type watchChan struct {
	watcher           *watcher
	key               string
	initialRev        int64
	recursive         bool
	ctx               context.Context
	cancel            context.CancelFunc
	incomingEventChan chan *event
	resultChan        chan watch.Event
	errChan           chan error
}

func newWatcher(
	collection *mongo.Collection,
	codec runtime.Codec,
	newFunc func() runtime.Object,
	versioner storage.Versioner,
) *watcher {
	res := &watcher{
		collection: collection,
		codec:      codec,
		newFunc:    newFunc,
		versioner:  versioner,
	}
	if newFunc == nil {
		res.objectType = "<unknown>"
	} else {
		res.objectType = reflect.TypeOf(newFunc()).String()
	}
	return res
}

func (w *watcher) Watch(ctx context.Context, key string, rev int64, recursive bool) (watch.Interface, error) {
	wc := w.createWatchChan(ctx, key, rev, recursive)
	go wc.run()
	return wc, nil
}

func (w *watcher) createWatchChan(ctx context.Context, key string, rev int64, recursive bool) *watchChan {
	wc := &watchChan{
		watcher:           w,
		key:               key,
		initialRev:        rev,
		recursive:         recursive,
		incomingEventChan: make(chan *event, incomingBufSize),
		resultChan:        make(chan watch.Event, outgoingBufSize),
		errChan:           make(chan error, 1),
	}

	wc.ctx, wc.cancel = context.WithCancel(ctx)
	return wc
}

func (wc *watchChan) run() {

	logrus.Tracef("starting watch on mongo collection %s/%s", wc.watcher.collection.Database().Name(), wc.watcher.collection.Name())

	watchClosedCh := make(chan struct{})
	go wc.startWatching(watchClosedCh)

	var resultChanWG sync.WaitGroup
	resultChanWG.Add(1)
	go wc.processEvent(&resultChanWG)

	select {
	case err := <-wc.errChan:
		if err == context.Canceled {
			break
		}
		errResult := transformErrorToEvent(err)
		if errResult != nil {
			select {
			case wc.resultChan <- *errResult:
			case <-wc.ctx.Done():
			}
		}
	case <-watchClosedCh:
	case <-wc.ctx.Done():
	}

	wc.cancel()

	resultChanWG.Wait()
	close(wc.resultChan)

}

func (wc *watchChan) processEvent(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case e := <-wc.incomingEventChan:
			res := wc.transform(e)
			if res == nil {
				continue
			}
			if len(wc.resultChan) == outgoingBufSize {
				logrus.Warn("fast watcher, slow processing")
			}
			select {
			case wc.resultChan <- *res:
			case <-wc.ctx.Done():
				return
			}
		case <-wc.ctx.Done():
			return
		}
	}
}

func (wc *watchChan) Stop() {
	wc.cancel()
}

func (wc *watchChan) ResultChan() <-chan watch.Event {
	return wc.resultChan
}

func (wc *watchChan) startWatching(watchClosedCh chan struct{}) {
	if wc.initialRev == 0 {
		if err := wc.sync(); err != nil {
			logrus.Errorf("failed to sync with latest state: %v", err)
			wc.sendError(err)
			return
		}
	}
	cs, err := wc.watcher.collection.Watch(wc.ctx, mongo.Pipeline{}, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if err != nil {
		logrus.Errorf("failed to start mongo change stream: %v", err)
		wc.sendError(err)
		return
	}

	for {
		if cs.Next(wc.ctx) {
			if cs.Err() != nil {
				err := cs.Err()
				wc.sendError(err)
				return
			}
			parsedEvent, err := parseEvent(cs.Current)
			if err != nil {
				logrus.Errorf("failed to parse mongo change stream event: %v", err)
				wc.sendError(err)
				return
			}
			if parsedEvent == nil {
				continue
			}
			wc.sendEvent(parsedEvent)
		} else {
			break
		}
	}
	close(watchClosedCh)
}

func (wc *watchChan) sendError(err error) {
	select {
	case wc.errChan <- err:
	case <-wc.ctx.Done():
	}
}

func (wc *watchChan) sendEvent(e *event) {
	if len(wc.incomingEventChan) == incomingBufSize {
		logrus.Warn("fast watcher, slow processing")
	}
	select {
	case wc.incomingEventChan <- e:
	case <-wc.ctx.Done():
	}
}

func parseDoc(versioner storage.Versioner, e bson.Raw) (*event, error) {

	key := e.Lookup("_id").ObjectID()

	current, err := e.LookupErr("current")
	if err != nil {
		return nil, err
	}

	currentVal := current.Value

	currentDocument, ok := current.DocumentOK()
	if !ok {
		// the "current" key was deleted, which means that the document was deleted.
		return nil, nil
	}

	metadata, err := currentDocument.LookupErr("metadata")
	if err != nil {
		return nil, err
	}

	metaDoc, ok := metadata.DocumentOK()
	if !ok {
		return nil, fmt.Errorf("cannot parse event: 'current.metadata' key is not a document")
	}

	resourceVersion, err := metaDoc.LookupErr("resourceVersion")
	if err != nil {
		return nil, err
	}

	resourceVersionStr, ok := resourceVersion.StringValueOK()
	if !ok {
		return nil, fmt.Errorf("cannot parse event: 'current.metadata.resourceVersion' key is not an int64")
	}

	resourceVersionInt, err := versioner.ParseResourceVersion(resourceVersionStr)
	if err != nil {
		return nil, err
	}

	ret := &event{
		key:       key.Hex(),
		value:     currentVal,
		prevValue: nil,
		rev:       resourceVersionInt,
		isDeleted: false,
		isCreated: true,
	}
	return ret, nil

}

func parseEvent(e bson.Raw) (*event, error) {

	eType := e.Lookup("operationType").StringValue()

	if eType == "delete" {
		return nil, nil
	}

	key, ok := e.Lookup("documentKey", "_id").ObjectIDOK()
	if !ok {
		return nil, fmt.Errorf("could not get document key")
	}

	previous, err := e.LookupErr("fullDocument", "previous")
	if err != nil {
		return nil, err
	}

	current, err := e.LookupErr("fullDocument", "current")
	if err != nil {
		return nil, err
	}

	currentVal := current.Value
	resourceVersion, err := e.LookupErr("fullDocument", "__revision")
	if err != nil {
		return nil, err
	}

	resourceVersionInt, ok := resourceVersion.AsInt64OK()
	if !ok {
		return nil, fmt.Errorf("cannot parse event: '__revision' key is not an int64")
	}

	ret := &event{
		key:       key.Hex(),
		value:     currentVal,
		prevValue: previous.Value,
		rev:       uint64(resourceVersionInt),
		isDeleted: len(currentVal) == 0 && len(previous.Value) > 0,
		isCreated: eType == "insert",
	}
	return ret, nil

}

func (wc *watchChan) sync() error {

	cursor, err := wc.watcher.collection.Find(wc.ctx, bson.M{})
	if err != nil {
		return err
	}
	for {
		if !cursor.Next(wc.ctx) {
			break
		}
		if cursor.Err() != nil {
			return cursor.Err()
		}
		doc, err := parseDoc(wc.watcher.versioner, cursor.Current)
		if err != nil {
			return err
		}
		if doc == nil {
			continue
		}
		wc.sendEvent(doc)
	}
	return nil
}

func transformErrorToEvent(err error) *watch.Event {
	//err = interpretWatchError(err)
	if _, ok := err.(exceptions.APIStatus); !ok {
		err = exceptions.NewInternalError(err)
	}
	status := err.(exceptions.APIStatus).Status()
	return &watch.Event{
		Type:   watch.Error,
		Object: &status,
	}
}

func (wc *watchChan) transform(e *event) (res *watch.Event) {

	curObj, oldObj, err := wc.prepareObjs(e)
	if err != nil {
		logrus.Errorf("failed to prepare current and previous objects: %v", err)
		wc.sendError(err)
		return nil
	}

	switch {
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
		res = &watch.Event{
			Type:   watch.Modified,
			Object: curObj,
		}
	}
	return res
}

func (wc *watchChan) prepareObjs(e *event) (curObj runtime.Object, oldObj runtime.Object, err error) {
	if !e.isDeleted {
		data, err := transformFromStorage(e.key, e.value)
		if err != nil {
			return nil, nil, err
		}
		curObj, err = decodeObj(wc.watcher.codec, wc.watcher.versioner, data, e.rev)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(e.prevValue) > 0 && (e.isDeleted) {
		data, err := transformFromStorage(e.key, e.prevValue)
		oldObj, err = decodeObj(wc.watcher.codec, wc.watcher.versioner, data, e.rev)
		if err != nil {
			return nil, nil, err
		}
	}
	return curObj, oldObj, nil
}

func decodeObj(codec runtime.Codec, versioner storage.Versioner, data []byte, rev uint64) (_ runtime.Object, err error) {
	obj, err := runtime.Decode(codec, data)
	if err != nil {
		return nil, err
	}

	if err := versioner.UpdateObject(obj, uint64(rev)); err != nil {
		return nil, fmt.Errorf("unable to version api object (%d) %#v: %v", rev, obj, err)
	}
	return obj, nil
}

func transformFromStorage(key string, data []byte) ([]byte, error) {

	var m = map[string]interface{}{}
	if err := bson.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	meta := m["metadata"].(map[string]interface{})
	meta["uid"] = key

	objBytes, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}

	return objBytes, nil
}
