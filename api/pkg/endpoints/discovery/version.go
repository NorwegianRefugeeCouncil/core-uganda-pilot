package discovery

import (
	"github.com/emicklei/go-restful"
	v1 "github.com/nrc-no/core/api/pkg/apis/discovery/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/negotiation"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
)

type APIResourceLister interface {
	ListAPIResources() []v1.APIResource
}

type APIResourceListerFunc func() []v1.APIResource

func (f APIResourceListerFunc) ListAPIResources() []v1.APIResource {
	return f()
}

type APIVersionHandler struct {
	serializer        runtime.NegotiatedSerializer
	groupVersion      schema.GroupVersion
	apiResourceLister APIResourceLister
}

func NewAPIVersionHandler(serializer runtime.NegotiatedSerializer, groupVersion schema.GroupVersion, apiResourceLister APIResourceLister) *APIVersionHandler {
	return &APIVersionHandler{
		serializer:        serializer,
		groupVersion:      groupVersion,
		apiResourceLister: apiResourceLister,
	}
}

func (s *APIVersionHandler) AddToWebService(ws *restful.WebService) {
	mediaTypes, _ := negotiation.MediaTypesForSerializer(s.serializer)
	ws.Route(ws.GET("/").To(s.handle).
		Doc("get available resources").
		Operation("getApiResources").
		Produces(mediaTypes...).
		Consumes(mediaTypes...).
		Writes(v1.APIResourceList{}))
}

func (h *APIVersionHandler) handle(req *restful.Request, resp *restful.Response) {
	h.ServeHTTP(resp.ResponseWriter, req.Request)
}

func (h *APIVersionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	responsewriters.WriteObjectNegotiated(
		h.serializer,
		negotiation.DefaultEndpointRestrictions,
		v1.SchemeGroupVersion,
		w,
		req, http.StatusOK,
		&v1.APIResourceList{
			GroupVersion: h.groupVersion.String(),
			APIResources: h.apiResourceLister.ListAPIResources(),
		},
	)
}
