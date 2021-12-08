package folder

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Get(folderId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		folder, err := h.store.Get(ctx, folderId)
		if err != nil {
			l.Error("failed to get folder from store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		utils.JSONResponse(w, http.StatusOK, folder)
	}
}

func (h *Handler) RestfulGet(request *restful.Request, response *restful.Response) {
	handler := h.Get(request.PathParameter(constants.ParamFolderID))
	handler(response.ResponseWriter, request.Request)
}
