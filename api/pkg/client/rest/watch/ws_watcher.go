package watch

import (
	"context"
	"fmt"
	v1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/watch"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

// WebSocketWatcher turns any stream for which you can write a Decoder interface
// into a watch.Interface.
type WebSocketWatcher struct {
	sync.Mutex
	ctx     context.Context
	source  func() ([]byte, error)
	stop    func()
	decoder runtime.Decoder
	result  chan watch.Event
	done    <-chan struct {
	}
}

// NewWebSocketWatcher creates a WebSocketWatcher from the given decoder.
func NewWebSocketWatcher(ctx context.Context, stop func(), d runtime.Decoder, source func() ([]byte, error)) *WebSocketWatcher {
	sw := &WebSocketWatcher{
		ctx:     ctx,
		decoder: d,
		source:  source,
		stop:    stop,
		// It's easy for a consumer to add buffering via an extra
		// goroutine/channel, but impossible for them to remove it,
		// so nonbuffered is better.
		result: make(chan watch.Event),
		// If the watcher is externally stopped there is no receiver anymore
		// and the send operations on the result channel, especially the
		// error reporting might block forever.
		// Therefore a dedicated stop channel is used to resolve this blocking.
		done: ctx.Done(),
	}
	go sw.receive()
	return sw
}

// ResultChan implements Interface.
func (sw *WebSocketWatcher) ResultChan() <-chan watch.Event {
	return sw.result
}

// Stop implements Interface.
func (sw *WebSocketWatcher) Stop() {
	// Call Close() exactly once by locking and setting a flag.
	sw.Lock()
	defer sw.Unlock()
	// closing a closed channel always panics, therefore check before closing
	select {
	case <-sw.done:
	default:
		sw.stop()
	}
}

// receive reads result from the decoder in a loop and sends down the result channel.
func (sw *WebSocketWatcher) receive() {
	defer close(sw.result)
	defer sw.Stop()
	for {

		data, err := sw.source()
		if err != nil {
			break
		}

		var evt = &v1.WatchEvent{}
		_, _, err = sw.decoder.Decode(data, nil, evt)
		if err != nil {
			break
		}

		fmt.Println("Received event raw " + string(evt.Object.Raw))

		obj, err := runtime.Decode(sw.decoder, evt.Object.Raw)
		if err != nil {
			logrus.Errorf("could not decode event object: %v", err)
			return
		}

		var watchEvent watch.Event
		if err := v1.Convert_v1_WatchEvent_To_watch_Event(evt, &watchEvent, nil); err != nil {
			break
		}

		if err != nil {
			switch err {
			case io.EOF:
				// watch closed normally
			case io.ErrUnexpectedEOF:
				logrus.Errorf("unexpected EOF during watch stream event decoding: %v", err)
			default:
				select {
				case <-sw.done:
				case sw.result <- watchEvent:
				}
			}
			return
		}
		select {
		case <-sw.done:
			return
		case sw.result <- watch.Event{
			Type:   watch.EventType(evt.Type),
			Object: obj,
		}:
		}
	}
}
