package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/meta"
	"github.com/nrc-no/core/pkg/logging"
	"github.com/nrc-no/core/pkg/server/core-db/types"
	"github.com/nrc-no/core/pkg/utils"
	"github.com/nrc-no/core/pkg/validation"
	"go.uber.org/zap"
)

func (c *Controller) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		l := logging.NewLogger(ctx)

		l.Debug("Unmarshaling entity")
		var entity types.Entity
		if err := utils.BindJSON(req, &entity); err != nil {
			l.Error("Failed to unmarshal entity", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("Validating entity")
		if errs := validation.Validate(&entity, true); !errs.IsEmpty() {
			l.Error("Failed to validate entity", zap.Error(errs.ToAggregate()))
			utils.ErrorResponse(w, meta.NewInvalid(types.EntityGR, "", errs))
			return
		}

		l.Debug("Storing entity")
		respEntity, err := c.entityService.Create(ctx, entity)
		if err != nil {
			l.Error("Failed to store entity", zap.Error(err))
			utils.ErrorResponse(w, err)
			return
		}

		l.Debug("Successfully created entity")
		utils.JSONResponse(w, http.StatusCreated, respEntity)
	}
}

func (c *Controller) RestfulCreate(request *restful.Request, response *restful.Response) {
	handler := c.Create()
	handler(response.ResponseWriter, request.Request)
}
