package organizations

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
)

type Handler struct {
	orgStore store.OrganizationStore
	ws       *restful.WebService
}

func NewHandler(orgStore store.OrganizationStore, idpStore store.IdentityProviderStore) (*Handler, error) {
	h := &Handler{}
	ws := new(restful.WebService).Path("/admin/organizations")
	ws.Route(ws.GET("").To(restfulList(orgStore)))
	ws.Route(ws.GET("/{organizationId}").To(restfulGet(orgStore, idpStore)))
	ws.Route(ws.POST("/{organizationId}").To(restfulUpdate(orgStore)))
	ws.Route(ws.GET("/add").To(restfulAddGet()))
	ws.Route(ws.POST("/add").To(restfulAddPost(orgStore)))
	h.ws = ws
	return h, nil
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}
