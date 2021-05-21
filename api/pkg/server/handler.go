package server

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"k8s.io/apiserver/pkg/server/mux"
	"net/http"
	"strings"
)

type HandlerChainBuilderFn func(apiHandler http.Handler) http.Handler

type APIServerHandler struct {
	GoRestfulContainer *restful.Container
	NonGoRestfulMux    *mux.PathRecorderMux
	FullHandlerChain   http.Handler
}

func NewAPIServerHandler(handlerChainBuilder HandlerChainBuilderFn) *APIServerHandler {
	goRestfulContainer := restful.NewContainer()
	nonGoRestfulMux := mux.NewPathRecorderMux("")
	nonGoRestfulMux.NotFoundHandler(http.NotFoundHandler())
	return &APIServerHandler{
		GoRestfulContainer: goRestfulContainer,
		NonGoRestfulMux:    nonGoRestfulMux,
		FullHandlerChain: handlerChainBuilder(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			// At this point, all the filters/middlewares have been ran
			// through the handlerChain.

			path := req.URL.Path

			// check to see if our webservices want to claim this path
			for _, ws := range goRestfulContainer.RegisteredWebServices() {
				switch {
				case ws.RootPath() == "/apis":
					// if we are exactly /apis or /apis/, then we need special handling in loop.
					// normally these are passed to the nonGoRestfulMux, but if discovery is enabled, it will go directly.
					// We can't rely on a prefix match since /apis matches everything (see the big comment on Director above)
					if path == "/apis" || path == "/apis/" {
						logrus.Infof("%v: %v %q satisfied by gorestful with webservice %v", "", req.Method, path, ws.RootPath())
						// don't use servemux here because gorestful servemuxes get messed up when removing webservices
						// TODO fix gorestful, remove TPRs, or stop using gorestful
						goRestfulContainer.Dispatch(w, req)
						return
					}

				case strings.HasPrefix(path, ws.RootPath()):
					// ensure an exact match or a path boundary match
					if len(path) == len(ws.RootPath()) || path[len(ws.RootPath())] == '/' {
						logrus.Infof("%v: %v %q satisfied by gorestful with webservice %v", "", req.Method, path, ws.RootPath())
						// don't use servemux here because gorestful servemuxes get messed up when removing webservices
						// TODO fix gorestful, remove TPRs, or stop using gorestful
						goRestfulContainer.Dispatch(w, req)
						return
					}
				}
			}

			// if we didn't find a match, then we just skip gorestful altogether
			logrus.Infof("%v: %v %q satisfied by nonGoRestful", "", req.Method, path)
			nonGoRestfulMux.ServeHTTP(w, req)

		})),
	}
}
