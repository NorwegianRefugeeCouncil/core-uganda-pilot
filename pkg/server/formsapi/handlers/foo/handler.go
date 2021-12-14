package foo

import (
	"net/http"
	"path"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/store"
)

type Handler struct {
	store      store.FooStore
	webService *restful.WebService
}

func NewHandler(store store.FooStore, rootPath string) *Handler {
	ws := new(restful.WebService).
		Path(path.Join(rootPath, "/foo")).
		Doc("foo.core.nrc.no API")

	h := &Handler{store: store, webService: ws}

	ws.Route(
		ws.POST("").
			To(h.RestfulCreate).
			Doc("Create a foo").
			Operation("createFoo").
			Consumes(mimetypes.ApplicationJson).
			Produces(mimetypes.ApplicationJson).
			Reads(types.FooReads{}).
			Writes(types.Foo{}).
			Returns(http.StatusCreated, "Created", types.Foo{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
