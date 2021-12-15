package folder

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.FolderStore
	webService *restful.WebService
}

func NewHandler(store store.FolderStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).
		Path("/apis/core.nrc.no/v1/folders").
		Doc("folders.core.nrc.no API")

	h.webService = ws

	folderPath := fmt.Sprintf("/{%s}", constants.ParamFolderID)

	ws.Route(ws.DELETE(folderPath).To(h.RestfulDelete).
		Doc("delete a folder").
		Operation("deleteFolder").
		Param(restful.PathParameter(constants.ParamFolderID, "id of the folder").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Returns(http.StatusNoContent, "OK", nil))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a folder").
		Operation("createFolder").
		Consumes(mimetypes.ApplicationJson).
		Produces(mimetypes.ApplicationJson).
		Reads(types.Folder{}).
		Writes(types.Folder{}).
		Returns(http.StatusOK, "OK", types.Folder{}))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("list all folders").
		Operation("listFolders").
		Produces(mimetypes.ApplicationJson).
		Writes(types.FolderList{}).
		Returns(http.StatusOK, "OK", types.FolderList{}))

	ws.Route(ws.GET(fmt.Sprintf("/{%s}", constants.ParamFolderID)).To(h.RestfulGet).
		Doc("get a folder").
		Operation("getFolder").
		Produces(mimetypes.ApplicationJson).
		Writes(types.Folder{}).
		Param(
			ws.
				PathParameter(constants.ParamFolderID, "id of the folder").
				Required(true).
				DataType("string").
				DataFormat("uuid"),
		).
		Returns(http.StatusOK, "OK", types.Folder{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
