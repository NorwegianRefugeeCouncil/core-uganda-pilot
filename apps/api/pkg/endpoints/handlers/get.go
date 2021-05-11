package handlers

import (
	"context"
	"encoding/json"
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

			w.Header().Set("Content-Type", serializer.Accepted.MediaType)
			u := websocket.Upgrader{}
			conn, err := u.Upgrade(w, req, nil)
			if err != nil {
				scope.Error(err, w, req)
				return
			}
			defer conn.Close()

			for {
				select {
				case <-ctx.Done():
					break
					return
				case evt := <-watcher.ResultChan():
					//if !ok {
					//  cancel()
					//  break
					//}

					evtBytes, err := json.Marshal(evt)
					if err != nil {
						break
					}
					logrus.Infof("sending json event: %s", string(evtBytes))

					if err := conn.WriteJSON(evt); err != nil {
						cancel()
						break
					}
				}
				break
			}
			return
		}

		result, err := r.List(ctx, &opts)
		if err != nil {
			scope.Error(err, w, req)
			return
		}

		transformResponseObject(ctx, scope, req, w, http.StatusOK, outputMediaType, result)

	}
}
