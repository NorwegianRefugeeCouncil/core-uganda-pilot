package folder

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/api/types/validation"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("unmarshaling folder")
		var folder types.Folder
		if err := utils.BindJSON(req, &folder); err != nil {
			l.Error("failed to unmarshal folder", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		folder.ID = uuid.NewV4().String()

		if errList := validation.ValidateFolder(&folder); errList.HasAny() {
			err := meta.NewInvalid(types.FolderGR, "", errList)
			l.Warn("folder is invalid", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("storing folder")
		respForm, err := h.store.Create(ctx, &folder)
		if err != nil {
			l.Error("failed to store folder", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("successfully created folder")
		utils.JSONResponse(w, http.StatusOK, respForm)
	}
}

func (h *Handler) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := h.Create()
	handler(response.ResponseWriter, request.Request)
}
