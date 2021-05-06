package server

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg3/endpoints"
	"github.com/nrc-no/core/apps/api/pkg3/endpoints/filters"
	"net/http"
)

type Server struct {
	Container    *restful.Container
	handlerChain http.Handler
}

func NewServer() *Server {
	var handlerChainBuilder HandlerChainBuilderFn = func(handler http.Handler) http.Handler {
		handler = filters.WithRequestInfo(handler, endpoints.DefaultRequestInfoFactory)
		return handler
	}
	goRestfulContainer := restful.NewContainer()
	apiServerHandler := NewAPIServerHandler(goRestfulContainer, handlerChainBuilder)

	server := &Server{
		Container:    goRestfulContainer,
		handlerChain: apiServerHandler,
	}

	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.handlerChain.ServeHTTP(w, req)
}

type HandlerChainBuilderFn func(apiHandler http.Handler) http.Handler
