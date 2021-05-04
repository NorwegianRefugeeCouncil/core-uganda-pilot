package server

import (
  "github.com/emicklei/go-restful"
  "github.com/sirupsen/logrus"
  "net/http"
  "strings"
)

type APIServerHandler struct {
  GoRestfulContainer *restful.Container
  Director           *director
  FullHandlerChain   http.Handler
}

type HandlerChainBuilderFn func(apiHandler http.Handler) http.Handler

func NewAPIServerHandler(name string, handlerChainBuilder HandlerChainBuilderFn) *APIServerHandler {

  goRestfulContainer := restful.NewContainer()
  goRestfulContainer.ServeMux = http.NewServeMux()
  goRestfulContainer.Router(restful.CurlyRouter{})

  director := &director{
    name:               name,
    goRestfulContainer: goRestfulContainer,
  }

  return &APIServerHandler{
    FullHandlerChain:   handlerChainBuilder(director),
    GoRestfulContainer: goRestfulContainer,
    Director:           director,
  }
}

func (a *APIServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
  a.FullHandlerChain.ServeHTTP(w, r)
}

type director struct {
  name               string
  goRestfulContainer *restful.Container
}

func (d *director) ServeHTTP(w http.ResponseWriter, req *http.Request) {

  path := req.URL.Path

  for _, ws := range d.goRestfulContainer.RegisteredWebServices() {
    switch {
    case ws.RootPath() == "/apis":
      if path == "/apis" || path == "/apis/" {
        logrus.Infof("%v: %v %q satisfied by gorestful with webservice %v", d.name, req.Method, path, ws.RootPath())
        d.goRestfulContainer.Dispatch(w, req)
        return
      }
    case strings.HasPrefix(path, ws.RootPath()):
      if len(path) == len(ws.RootPath()) || string(path[len(ws.RootPath())]) == "/" {
        logrus.Infof("%v: %v %q satisfied by gorestful with webservice %v", d.name, req.Method, path, ws.RootPath())
        d.goRestfulContainer.Dispatch(w, req)
        return
      }
    }
  }

  http.Error(w, "not found", http.StatusNotFound)

}
