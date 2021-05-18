package server

import (
  "github.com/emicklei/go-restful"
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/server/mux"
  "k8s.io/klog/v2"
  "net/http"
  "strings"
)

type APIServerHandler struct {
  GoRestfulContainer *restful.Container
  HandlerChain       http.Handler
  NonGoRestfulMux    *mux.PathRecorderMux
}

func NewAPIServerHandler(
  name string,
  s runtime.NegotiatedSerializer,
  handlerChainBuilder HandlerChainBuilderFn,
  notFoundHandler http.Handler,
) *APIServerHandler {

  nonGoRestfulMux := mux.NewPathRecorderMux("")
  if notFoundHandler != nil {
    nonGoRestfulMux.NotFoundHandler(notFoundHandler)
  }

  goRestfulContainer := restful.NewContainer()
  handler := &APIServerHandler{
    GoRestfulContainer: goRestfulContainer,
    HandlerChain: handlerChainBuilder(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

      path := req.URL.Path
      // check to see if our webservices want to claim this path
      for _, ws := range goRestfulContainer.RegisteredWebServices() {
        switch {
        case ws.RootPath() == "/apis":
          // if we are exactly /apis or /apis/, then we need special handling in loop.
          // normally these are passed to the nonGoRestfulMux, but if discovery is enabled, it will go directly.
          // We can't rely on a prefix match since /apis matches everything (see the big comment on Director above)
          if path == "/apis" || path == "/apis/" {
            klog.V(5).Infof("%v: %v %q satisfied by gorestful with webservice %v", name, req.Method, path, ws.RootPath())
            // don't use servemux here because gorestful servemuxes get messed up when removing webservices
            // TODO fix gorestful, remove TPRs, or stop using gorestful
            goRestfulContainer.Dispatch(w, req)
            return
          }

        case strings.HasPrefix(path, ws.RootPath()):
          // ensure an exact match or a path boundary match
          if len(path) == len(ws.RootPath()) || path[len(ws.RootPath())] == '/' {
            klog.V(5).Infof("%v: %v %q satisfied by gorestful with webservice %v", name, req.Method, path, ws.RootPath())
            // don't use servemux here because gorestful servemuxes get messed up when removing webservices
            // TODO fix gorestful, remove TPRs, or stop using gorestful
            goRestfulContainer.Dispatch(w, req)
            return
          }
        }
      }

      // if we didn't find a match, then we just skip gorestful altogether
      klog.V(5).Infof("%v: %v %q satisfied by nonGoRestful", name, req.Method, path)
      nonGoRestfulMux.ServeHTTP(w, req)

      goRestfulContainer.Dispatch(w, req)
    })),
  }
  return handler
}

func (h *APIServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  h.HandlerChain.ServeHTTP(w, req)
}
