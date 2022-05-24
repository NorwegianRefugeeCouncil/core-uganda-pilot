package identity

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/api/mimetypes"
	"github.com/nrc-no/core/pkg/api/types"
	"github.com/nrc-no/core/pkg/constants"
	"github.com/nrc-no/core/pkg/store"
	"net/http"
)

type Handler struct {
	store      store.IdentityProfileStore
	webService *restful.WebService
}

func (h *Handler) WebService() *restful.WebService {
	return h.webService
}

func NewHandler(store store.IdentityProfileStore) *Handler {
	h := &Handler{store: store}

	ws := new(restful.WebService).
		Path("/apis/admin.nrc.no/v1/identities").
		Doc("identities.admin.nrc.no API")

	h.webService = ws

	ws.Route(ws.GET(fmt.Sprintf("/{%s}", constants.ParamIdentityID)).To(h.RestfulGet).
		Doc("gets an identity").
		Operation("getIdentity").
		Param(restful.PathParameter(constants.ParamIdentityID, "identity id").
			DataType("string").
			DataFormat("uuid").
			Required(true)).
		Produces(mimetypes.ApplicationJson).
		Writes(types.IdentityProfile{}).
		Returns(http.StatusOK, "OK", types.IdentityProfile{}))

	return h
}
