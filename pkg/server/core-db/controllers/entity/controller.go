package entity

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/common"
	model "github.com/nrc-no/core/pkg/server/core-db/models/entity"
	service "github.com/nrc-no/core/pkg/server/core-db/services/entity"
	"github.com/nrc-no/core/pkg/server/core-db/types"
)

type Controller struct {
	entityService      service.EntityService
	entityModel        model.EntityModel
	transactionManager common.TransactionManager
	webService         *restful.WebService
}

func NewController(
	entityModel model.EntityModel,
	transactionManager common.TransactionManager,
) *Controller {
	c := &Controller{
		entityService:      service.NewEntityService(entityModel, transactionManager),
		entityModel:        entityModel,
		transactionManager: transactionManager,
	}

	ws := new(restful.WebService).
		Path("/apis/core.nrc.no/v1/entity").
		Doc("entity.core.nrc.no API")

	c.webService = ws

	ws.Route(
		ws.
			POST("/").
			To(c.RestfulCreate).
			Doc("Create an entity").
			Operation("createEntity").
			Consumes(mimetypes.ApplicationJson).
			Produces(mimetypes.ApplicationJson).
			Reads(types.Entity{}).
			Writes(types.Entity{}).
			Returns(http.StatusOK, "OK", types.Entity{}),
	)

	ws.Route(
		ws.
			GET("/{entityID}").
			To(c.RestfulGet).
			Doc("Get an entity").
			Operation("getEntity").
			Param(ws.PathParameter("entityID", "Entity ID").DataType("string")).
			Produces(mimetypes.ApplicationJson).
			Writes(types.Entity{}).
			Returns(http.StatusOK, "OK", types.Entity{}),
	)

	ws.Route(
		ws.
			GET("/").
			To(c.RestfulList).
			Doc("List entities").
			Operation("listEntities").
			Produces(mimetypes.ApplicationJson).
			Writes([]types.Entity{}).
			Returns(http.StatusOK, "OK", []types.Entity{}),
	)

	ws.Route(
		ws.
			PUT("/{entityID}").
			To(c.RestfulUpdate).
			Doc("Update an entity").
			Operation("updateEntity").
			Param(ws.PathParameter("entityID", "Entity ID").DataType("string")).
			Consumes(mimetypes.ApplicationJson).
			Produces(mimetypes.ApplicationJson).
			Reads(types.Entity{}).
			Writes(types.Entity{}).
			Returns(http.StatusOK, "OK", types.Entity{}),
	)

	ws.Route(
		ws.
			DELETE("/{entityID}").
			To(c.RestfulDelete).
			Doc("Delete an entity").
			Operation("deleteEntity").
			Param(ws.PathParameter("entityID", "Entity ID").DataType("string")).
			Produces(mimetypes.ApplicationJson).
			Returns(http.StatusOK, "OK", types.Entity{}),
	)

	return c
}

func (c *Controller) WebService() *restful.WebService {
	return c.webService
}
