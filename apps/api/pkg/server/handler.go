package server

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"net/http"
)

type APIServerHandler struct {
	GoRestfulContainer *restful.Container
	HandlerChain       http.Handler
}

func NewAPIServerHandler(
	s runtime.NegotiatedSerializer,
	handlerChainBuilder HandlerChainBuilderFn,
	notFoundHandler http.Handler,
) *APIServerHandler {
	goRestfulContainer := restful.NewContainer()
	handler := &APIServerHandler{
		GoRestfulContainer: goRestfulContainer,
		HandlerChain: handlerChainBuilder(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			goRestfulContainer.Dispatch(writer, request)
		})),
	}
	return handler
}

func (h *APIServerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.HandlerChain.ServeHTTP(w, req)
}
