package server

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkgs2/endpoints"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Server struct {
	GoRestfulContainer *restful.Container
	Handler            http.Handler
}

func NewServer(handler *Handler) *Server {
	return &Server{
		Handler: handler,
	}
}

type Handler struct {
	goRestfulContainer *restful.Container
}

func NewHandler(goRestfulContainer *restful.Container) *Handler {
	return &Handler{goRestfulContainer: goRestfulContainer}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	for _, ws := range h.goRestfulContainer.RegisteredWebServices() {
		if strings.HasPrefix(ws.RootPath(), path) {
			h.goRestfulContainer.Dispatch(w, req)
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
	return
}

func (s *Server) installAPIResources(apiGroupInfo *endpoints.APIGroupInfo) error {
	for _, groupVersion := range apiGroupInfo.PrioritizedVersions {
		if len(apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version]) == 0 {
			logrus.Warn("skipping API %s because it has no resources", groupVersion)
			continue
		}

		apiGroupVersion := s.getAPIGroupVersion(apiGroupInfo, groupVersion)
		resourceInfo, err := apiGroupVersion.InstallREST(s.Handler, s.GoRestfulContainer)
		if err != nil {
			return fmt.Errorf("unable to setup API %v: %v", apiGroupInfo, err)
		}

	}
}

func (s *Server) getAPIGroupVersion(apiGroupInfo *endpoints.APIGroupInfo, groupVersion schema.GroupVersion) *endpoints.APIGroupVersion {
	storage := make(map[string]endpoints.Storage)
	for k, v := range apiGroupInfo.VersionedResourcesStorageMap {
		storage[strings.ToLower(k)] = v
	}
	version := s.newAPIGroupVersion(apiGroupInfo, groupVersion)
	version.Storage = storage
	return version
}

func (s *Server) newAPIGroupVersion(apiGroupInfo *endpoints.APIGroupInfo, groupVersion schema.GroupVersion) *endpoints.APIGroupVersion {
	return &endpoints.APIGroupVersion{
		Storage:      nil,
		GroupVersion: groupVersion,
		Typer:        apiGroupInfo.Scheme,
		Creater:      apiGroupInfo.Scheme,
	}
}
