package handlers

import (
	"bytes"
	"context"
	"github.com/gorilla/websocket"
	v1 "github.com/nrc-no/core/api/pkg/apis/meta/v1"
	rest2 "github.com/nrc-no/core/api/pkg/registry/rest"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"math/rand"
	"net/http"
	"time"
)

// ListResource is a generic REST handler to list resources (multiple result)
func ListResource(scope *RequestScope, getter rest2.Lister, watcher rest2.Watcher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		// Negotiate the output media type
		outputMediaType, serializerInfo, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		// Retrieve the v1.ListResourcesOptions
		opts := v1.ListResourcesOptions{}
		if err := v1.ParameterCodec.DecodeParameters(req.URL.Query(), v1.SchemeGroupVersion, &opts); err != nil {
			err = errors.NewBadRequest(err.Error())
			scope.err(err, w, req)
			return
		}

		// If this is a WATCH request, delegate to serveWatch
		if opts.Watch {
			serveWatch(ctx, w, req, watcher, scope, opts, outputMediaType, serializerInfo)
			return
		}

		// At this point, this is not a WATCH request
		result, err := getter.List(ctx)
		if err != nil {
			scope.err(err, w, req)
			return
		}

		transformResponseObject(scope, req, w, http.StatusOK, outputMediaType, result)
	}
}

// serveWatch will
// - upgrade the connection to a WebSocket connection
// - listen to the rest2.Watcher events
// - encode the events
// - send them through the websocket connection
func serveWatch(
	ctx context.Context,
	w http.ResponseWriter,
	req *http.Request,
	watcher rest2.Watcher,
	scope *RequestScope,
	opts v1.ListResourcesOptions,
	outputMediaType negotiation.MediaTypeOptions,
	serializerInfo runtime.SerializerInfo,
) {

	// Perhaps the REST endpoint does not support the rest.Watcher interface
	// In this case, throw a MethodNotAllowed error (HTTP 405)
	if watcher == nil {
		scope.err(errors.NewMethodNotSupported(scope.Resource.GroupResource(), "watch"), w, req)
		return
	}

	// We figure out the watch timeout. It's either provided in the ListResourcesOptions
	// Or we default to a random value between 30 and 60 minutes
	timeout := time.Duration(0)
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	if timeout == 0 {
		timeout = time.Duration(float64(30*time.Minute) * (rand.Float64() + 1.0))
	}

	//ctx, cancel := context.WithTimeout(ctx, timeout)
	//defer cancel()

	// Start watching
	wc, err := watcher.Watch(ctx, opts)
	if err != nil {
		scope.err(err, w, req)
		return
	}
	defer wc.Stop()

	// Upgrade the connection to a websocket connection
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, req, http.Header{
		"Content-Type": []string{
			outputMediaType.Accepted.MediaType,
		},
	})
	if err != nil {
		scope.err(err, w, req)
		return
	}
	defer ws.Close()

	kind, serializer, _ := targetEncodingForTransform(scope, outputMediaType, req)

	_, outSerializerInfo, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
	encoder := serializer.EncoderForVersion(outSerializerInfo.Serializer, kind.GroupVersion())

	// Wait and process watch events
	for {
		select {
		case evt := <-wc.ResultChan():

			if evt.Object == nil && len(evt.Type) == 0 {
				return
			}

			// Encode the watch.Event object
			buf := &bytes.Buffer{}
			if err := encoder.Encode(evt.Object, buf); err != nil {
				logrus.Errorf("unable to encode watch object %#v: %v", evt.Object, err)
				return
			}

			// Create a runtime.Unknown object to hold the encoded data
			uk := &runtime.Unknown{}
			uk.Raw = buf.Bytes()
			evt.Object = uk

			// Convert the watch.Event object to metav1.WatchEvent object
			outEvent := &metav1.WatchEvent{}
			if err := metav1.Convert_watch_Event_To_v1_WatchEvent(&evt, outEvent, nil); err != nil {
				logrus.Errorf("unable to convert to metav1.WatchEvent: %v", err)
				return
			}

			// Encode the metav1.WatchEvent
			outBytes, err := runtime.Encode(serializerInfo.Serializer, outEvent)
			if err != nil {
				logrus.Errorf("unable to encode metav1.WatchEvent: %v", err)
				return
			}

			// Send the encoded metav1.WatchEvent on the wire
			if err := ws.WriteMessage(websocket.TextMessage, outBytes); err != nil {
				logrus.Errorf("unable to write message: %v", err)
				return
			}

		case <-ctx.Done():
			return
		}
	}
	return
}
