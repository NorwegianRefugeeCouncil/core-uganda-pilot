package server

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type APIServerHandler struct {
	GoRestfulContainer *restful.Container
	HandlerChain       http.Handler
}

func NewAPIServerHandler(container *restful.Container, handlerChainBuilder HandlerChainBuilderFn) *APIServerHandler {
	handler := &APIServerHandler{
		GoRestfulContainer: container,
	}

	handlerChain := handlerChainBuilder(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.GoRestfulContainer.Dispatch(writer, request)
	}))
	handler.HandlerChain = handlerChain

	return handler
}

func (h *APIServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.HandlerChain.ServeHTTP(w, req)
}
