package login

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/nrc-no/core/pkg/store"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"net/url"
)

type Handler struct {
	hydraAdmin admin.ClientService
	ws         *restful.WebService
}

func NewHandler(orgStore store.OrganizationStore, idpStore store.IdentityProviderStore) (*Handler, error) {
	h := &Handler{}
	adminURL, err := url.Parse("http://localhost:4445")
	if err != nil {
		return nil, err
	}
	hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}).Admin
	h.hydraAdmin = hydraAdmin

	ws := new(restful.WebService)
	ws.Route(ws.GET("/login").To(restfulGetSubject(hydraAdmin)))
	ws.Route(ws.POST("/login").To(restfulGetLogin(hydraAdmin, orgStore, idpStore)))
	h.ws = ws
	return h, nil
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}
