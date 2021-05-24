package discovery

import (
	"github.com/emicklei/go-restful"
	v1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
)

type APIGroupHandler struct {
	serializer runtime.NegotiatedSerializer
	group      v1.APIGroup
}

func NewAPIGroupHandler(serializer runtime.NegotiatedSerializer, group v1.APIGroup) *APIGroupHandler {
	return &APIGroupHandler{
		serializer: serializer,
		group:      group,
	}
}

func (h *APIGroupHandler) WebService() *restful.WebService {
	mediaTypes, _ := negotiation.MediaTypesForSerializer(h.serializer)
	ws := new(restful.WebService)
	ws.Path("/apis/" + h.group.Name)
	ws.Doc("get information on API group " + h.group.Name)
	ws.Route(ws.GET("/").To(h.handle).
		Doc("get information on API group " + h.group.Name).
		Operation("getApiGroup").
		Produces(mediaTypes...).
		Consumes(mediaTypes...).
		Writes(v1.APIGroup{}))
	return ws
}

func (h *APIGroupHandler) handle(req *restful.Request, resp *restful.Response) {
	h.ServeHTTP(resp.ResponseWriter, req.Request)
}

func (h *APIGroupHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	responsewriters.WriteObjectNegotiated(
		h.serializer,
		negotiation.DefaultEndpointRestrictions,
		v1.SchemeGroupVersion,
		w,
		req,
		http.StatusOK,
		&h.group)
}
