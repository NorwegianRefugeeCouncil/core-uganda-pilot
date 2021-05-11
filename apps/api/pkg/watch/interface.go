package watch

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/sirupsen/logrus"
	"io"
)

// EventType defines the possible types of events.
type EventType string

const (
	Added    EventType = "ADDED"
	Modified EventType = "MODIFIED"
	Deleted  EventType = "DELETED"
	Bookmark EventType = "BOOKMARK"
	Error    EventType = "ERROR"
)

type Event struct {
	Type   EventType      `json:"type"`
	Object runtime.Object `json:"object"`
}

type Interface interface {
	Stop()
	ResultChan() <-chan Event
}

type watcher struct {
	resultChan       chan Event
	resultChanClosed bool
	errChan          chan error
	ctx              context.Context
	cancel           context.CancelFunc
	readNext         func() ([]byte, error)
	decoder          runtime.Decoder
}

func NewWatcher(ctx context.Context, readNext func() ([]byte, error), decoder runtime.Decoder) *watcher {
	watcher := &watcher{
		resultChan: make(chan Event, 100),
		errChan:    make(chan error, 1),
		readNext:   readNext,
		decoder:    decoder,
	}
	watcher.ctx, watcher.cancel = context.WithCancel(ctx)
	go watcher.run()
	return watcher
}

func (w *watcher) Stop() {
	w.cancel()
}

func (w *watcher) run() {
	go func() {
		for {
			payload, err := w.readNext()
			if err != nil {
				w.errChan <- err
				break
			}

			type tmp struct {
				Type   EventType              `json:"type"`
				Object map[string]interface{} `json:"object"`
			}
			var event tmp

			if err := json.Unmarshal(payload, &event); err != nil {
				w.errChan <- err
				break
			}

			runtimeObjectBytes, err := json.Marshal(event.Object)
			if err != nil {
				w.errChan <- err
				break
			}

			logrus.Infof("received event with type '%s' and content: %s\n", event.Type, string(runtimeObjectBytes))

			obj, _, err := w.decoder.Decode(runtimeObjectBytes, nil, nil)
			if err != nil {
				logrus.Errorf("unable to decode event data: %v", err)
				w.errChan <- err
				break
			}

			if w.resultChanClosed {
				break
			}

			w.resultChan <- Event{
				Type:   event.Type,
				Object: obj,
			}

		}
	}()

	select {
	case err := <-w.errChan:
		if errors.Is(err, io.EOF) {
			break
		}
		if err == context.Canceled {
			break
		}
	case <-w.ctx.Done():
	}

	w.cancel()
	close(w.resultChan)
	w.resultChanClosed = true

}

func (w *watcher) ResultChan() <-chan Event {
	return w.resultChan
}
