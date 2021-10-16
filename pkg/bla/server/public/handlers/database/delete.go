package database

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) Delete(databaseId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		_, err := uuid.FromString(databaseId)
		if err != nil {
			utils.ErrorResponse(w, meta.NewBadRequest("invalid database id"))
			return
		}

		if err := h.store.Delete(req.Context(), databaseId); err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *Handler) RestfulDelete(request *restful.Request, response *restful.Response) {
	databaseID := request.PathParameter(constants.ParamDatabaseID)
	h.Delete(databaseID)(response.ResponseWriter, request.Request)
}
