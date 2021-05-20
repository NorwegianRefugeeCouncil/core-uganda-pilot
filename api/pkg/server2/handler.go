package server2

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type HandlerChainBuilderFn func(apiHandler http.Handler) http.Handler

type APIServerHandler struct {
	GoRestfulContainer *restful.Container
	FullHandlerChain   http.Handler
}

func NewAPIServerHandler(handlerChainBuilder HandlerChainBuilderFn) *APIServerHandler {
	goRestfulContainer := restful.NewContainer()
	return &APIServerHandler{
		GoRestfulContainer: goRestfulContainer,
		FullHandlerChain: handlerChainBuilder(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			goRestfulContainer.Dispatch(writer, request)
		})),
	}
}
