package organization

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/bla/constants"
	"github.com/nrc-no/core/pkg/utils"
	"net/http"
)

func (h *Handler) Get(organizationID string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		result, err := h.store.Get(req.Context(), organizationID)
		if err != nil {
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, result)

	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	organizationID := request.PathParameter(constants.ParamOrganizationID)
	handler := h.Get(organizationID)
	handler(response.ResponseWriter, request.Request)
}
