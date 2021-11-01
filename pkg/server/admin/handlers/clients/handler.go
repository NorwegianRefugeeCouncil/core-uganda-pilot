package clients

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"net/http"
	"net/url"
)

type Handler struct {
	hydraAdmin admin.ClientService
	ws         *restful.WebService
}

func NewHandler() (*Handler, error) {
	h := &Handler{}
	adminURL, err := url.Parse("http://localhost:4445")
	if err != nil {
		return nil, err
	}
	hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path}).Admin
	h.hydraAdmin = hydraAdmin

	ws := new(restful.WebService).Path("/admin/clients/")

	ws.Route(ws.PUT("/{clientId}").To(restfulUpdate(hydraAdmin)).
		Doc(`updates oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Operation("updateClient").
		Reads(&Client{}).
		Writes(&Client{}).
		Returns(http.StatusOK, "OK", &Client{}),
	)

	ws.Route(ws.POST("/").To(restfulCreate(hydraAdmin)).
		Doc(`creates oauth2 client`).
		Operation("createClient").
		Reads(&Client{}).
		Writes(&Client{}).
		Returns(http.StatusOK, "OK", &Client{}),
	)

	ws.Route(ws.GET("/{clientId}").To(restfulGet(hydraAdmin)).
		Doc(`gets oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Operation("getClient").
		Reads(nil).
		Writes(&Client{}).
		Returns(http.StatusOK, "OK", &Client{}),
	)

	ws.Route(ws.DELETE("/{clientId}").To(restfulDelete(hydraAdmin)).
		Doc(`deletes oauth2 client`).
		Param(ws.PathParameter("clientId", "client id").Required(true)).
		Operation("deleteClient").
		Reads(nil).
		Writes(nil).
		Returns(http.StatusOK, "OK", &Client{}),
	)

	ws.Route(ws.GET("/").To(restfulList(hydraAdmin)).
		Doc(`gets oauth2 clients`).
		Operation("listClients").
		Reads(nil).
		Writes(&ClientList{}).
		Returns(http.StatusOK, "OK", &ClientList{}),
	)

	h.ws = ws
	return h, nil
}

func (h *Handler) WebService() *restful.WebService {
	return h.ws
}
