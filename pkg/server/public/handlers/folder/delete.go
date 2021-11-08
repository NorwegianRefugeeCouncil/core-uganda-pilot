package folder

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Delete(folderId string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx).With(zap.String("folder_id", folderId))

		l.Debug("validating folder id")
		_, err := uuid.FromString(folderId)
		if err != nil {
			l.Error("failed to validate folder id", zap.Error(err))
			utils.ErrorResponse(w, meta.NewBadRequest("invalid folder id"))
			return
		}

		l.Debug("deleting folder from store")
		if err := h.store.Delete(ctx, folderId); err != nil {
			l.Error("failed to delete folder in store", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully deleted folder")
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func (h *Handler) RestfulDelete(request *restful.Request, response *restful.Response) {
	h.Delete(request.PathParameter(constants.ParamFolderID))(response.ResponseWriter, request.Request)
}
