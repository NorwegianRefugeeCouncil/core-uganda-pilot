package main

import (
	"context"
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkgs2/client/nrc"
	"github.com/nrc-no/core/apps/api/pkgs2/client/rest"
	"github.com/nrc-no/core/apps/api/pkgs2/endpoints"
	server2 "github.com/nrc-no/core/apps/api/pkgs2/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	apiServer  *server2.Server
	httpServer *httptest.Server
	httpClient *http.Client
	restClient *rest.RESTClient
	nrcClient  *nrc.NrcCoreClient
)

func init() {
	goRestfulContainer := restful.NewContainer()
	goRestfulContainer.Router(restful.CurlyRouter{})

	endpoints.Register("core.nrc.no", "v1", "formdefinitions", goRestfulContainer)

	handler := server2.NewHandler(goRestfulContainer)
	apiServer = server2.NewServer(handler)
	httpServer = httptest.NewServer(apiServer.Handler)
	httpClient = httpServer.Client()
	serverURL, err := url.Parse(httpServer.URL)
	if err != nil {
		panic(err)
	}
	restClient = rest.NewRESTClient(serverURL, "/apis/core.nrc.no/v1", rest.RESTClientConfig{}, httpClient)
	nrcClient = nrc.New(restClient)
}

func TestServer(t *testing.T) {
	ctx := context.TODO()
	formDefinitions, err := nrcClient.FormDefinitions().List(ctx)
	if !assert.NoError(t, err) {
		return
	}

	t.Logf("%#v", formDefinitions)

}
