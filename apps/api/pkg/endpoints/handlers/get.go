package handlers

import (
	"bytes"
	"context"
	"github.com/gorilla/websocket"
	"github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion"
	metainternalscheme "github.com/nrc-no/core/apps/api/pkg/apis/meta/internalversion/scheme"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/negotiation"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"github.com/sirupsen/logrus"
	"net/http"
)

type getterFunc func(ctx context.Context, name string, req *http.Request) (runtime.Object, error)

func getResourceHandler(scope *RequestScope, getter getterFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, name, err := scope.Namer.Name(req)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		ctx := req.Context()
		result, err := getter(ctx, name, req)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

	}
}

func GetResource(r rest.Getter, scope *RequestScope) http.HandlerFunc {
	return getResourceHandler(scope, func(ctx context.Context, name string, req *http.Request) (runtime.Object, error) {
		options := metav1.GetOptions{}
		if values := req.URL.Query(); len(values) > 0 {
			if err := metainternalscheme.ParameterCodec.DecodeParameters(values, metav1.SchemeGroupVersion, &options); err != nil {
				err = exceptions.NewBadRequest(err.Error())
				return nil, err
			}
		}
		return r.Get(ctx, name, &options)
	})
}

func ListResource(r rest.Lister, rw rest.Watcher, scope *RequestScope) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		outputMediaType, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		serializer, _, err := negotiation.NegotiateOutputMediaType(req, scope.Serializer, scope)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		opts := internalversion.ListOptions{}
		if err := metainternalscheme.ParameterCodec.DecodeParameters(req.URL.Query(), metav1.SchemeGroupVersion, &opts); err != nil {
			err = exceptions.NewBadRequest(err.Error())
			scope.Error(err, w, req)
			return
		}

		if opts.Watch {

			if rw == nil {
				scope.Error(exceptions.NewMethodNotSupported(scope.Resource.GroupResource(), "watch"), w, req)
				return
			}

			logrus.Infof("starting watch path: %s, resourceVersion %s", req.URL.String(), scope.Resource.Resource)
			watcher, err := rw.Watch(ctx, &opts)
			if err != nil {
				scope.Error(err, w, req)
				return
			}

			responseHeader := http.Header{}
			responseHeader.Set("Content-Type", serializer.Accepted.MediaType)
			u := websocket.Upgrader{}
			conn, err := u.Upgrade(w, req, responseHeader)
			if err != nil {
				scope.Error(err, w, req)
				return
			}
			defer conn.Close()

			for {
				select {
				case <-ctx.Done():
					return
				case evt := <-watcher.ResultChan():
					//if !ok {
					//  cancel()
					//  break
					//}

					encoder := scope.Serializer.EncoderForVersion(serializer.Accepted.Serializer, scope.Kind.GroupVersion())
					embeddedEncoder := scope.Serializer.EncoderForVersion(serializer.Accepted.Serializer, scope.Kind.GroupVersion())

					var unknown runtime.Unknown
					internalEvent := &metav1.InternalEvent{}
					buf := bytes.NewBuffer(nil)
					if err := embeddedEncoder.Encode(evt.Object, buf); err != nil {
						scope.Error(err, w, req)
						return
					}
					unknown.Raw = buf.Bytes()
					evt.Object = &unknown

					outEvent := &metav1.WatchEvent{}
					*internalEvent = metav1.InternalEvent(evt)

					err = metav1.Convert_v1_InternalEvent_To_v1_WatchEvent(internalEvent, outEvent, nil)
					if err != nil {
						logrus.Errorf("unable to convert internal event to watch event")
					}

					var outBuf = bytes.NewBuffer(nil)
					if err := encoder.Encode(outEvent, outBuf); err != nil {
						logrus.Errorf("unable to serialize out event: %v", err)
						return
					}

					if err := conn.WriteMessage(websocket.TextMessage, outBuf.Bytes()); err != nil {
						logrus.Errorf("failed to write to websocket: %v", err)
						return
					}

				}
			}
		}

		result, err := r.List(ctx, &opts)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

	}
}
