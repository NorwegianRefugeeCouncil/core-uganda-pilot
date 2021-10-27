package folder

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) Delete(folderId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, err := uuid.FromString(folderId)
		if err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid folder id"))
			return
		}

		if err := h.store.Delete(req.Context(), folderId); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *Handler) RestfulDelete(request *restful.Request, response *restful.Response) {
	h.Delete(request.PathParameter(constants.ParamFolderID))(response.ResponseWriter, request.Request)
}
