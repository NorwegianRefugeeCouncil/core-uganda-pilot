package form

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Delete(formId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("deleting form")
		if err := h.store.Delete(ctx, formId); err != nil {
			l.Error("failed to delete form", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully deleted form")
		utils.JSONResponse(w, http.StatusNoContent, nil)
	}
}

func (h *Handler) RestfulDelete(request *restful.Request, response *restful.Response) {
	handler := h.Delete(request.PathParameter(constants.ParamFormID))
	handler(response.ResponseWriter, request.Request)
}
