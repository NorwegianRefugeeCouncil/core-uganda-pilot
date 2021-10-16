package folder

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/bla/store"
	"github.com/nrc-no/core/pkg/bla/types"
	"net/http"
)

type Handler struct {
	store      store.FolderStore
	webService *restful.WebService
}

func NewHandler(store store.FolderStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).Path("/folders")
	h.webService = ws

	folderPath := fmt.Sprintf("/{%s}", constants.ParamFolderID)

	ws.Route(ws.DELETE(folderPath).To(h.RestfulDelete).
		Doc("delete a folder").
		Operation("deleteFolder").
		Param(restful.PathParameter(constants.ParamFolderID, "id of the folder").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Writes(nil).
		Returns(http.StatusNoContent, "OK", nil))

	ws.Route(ws.POST("/").To(h.RestfulCreate).
		Doc("create a folder").
		Operation("createFolder").
		Reads(types.Folder{}).
		Writes(types.Folder{}).
		Returns(http.StatusOK, "OK", types.Folder{}))

	ws.Route(ws.GET("/").To(h.RestfulList).
		Doc("list all folders").
		Operation("listFolders").
		Writes(types.FolderList{}).
		Returns(http.StatusOK, "OK", types.FolderList{}))

	return h
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}
