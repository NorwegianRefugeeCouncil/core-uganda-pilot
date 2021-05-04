package server

import (
  "fmt"
  "github.com/nrc-no/core/apps/api/pkg/endpoints"
  "github.com/nrc-no/core/apps/api/pkg/registry/rest"
  "github.com/sirupsen/logrus"
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/apimachinery/pkg/util/sets"
  "strings"
)

type Server struct {
  Handler                *APIServerHandler
  legacyAPIGroupPrefixes sets.String
}

type APIGroupInfo struct {
  PrioritizedVersions          []schema.GroupVersion
  VersionedResourcesStorageMap map[string]map[string]rest.Storage
  MetaGroupVersion             *schema.GroupVersion
  ParameterCodec               runtime.ParameterCodec
  NegotiatedSerializer         runtime.NegotiatedSerializer
  Scheme                       *runtime.Scheme
}

func (s *Server) InstallAPIGroups(apiGroupInfos ...*APIGroupInfo) error {

  for _, info := range apiGroupInfos {
    if len(info.PrioritizedVersions[0].Group) == 0 {
      return fmt.Errorf("cannot register handler with empty group for %#v", *info)
    }
    if len(info.PrioritizedVersions[0].Version) == 0 {
      return fmt.Errorf("cannot register handler with an empty version for %#v", *info)
    }
  }

  for _, apiGroupInfo := range apiGroupInfos {
    if err := s.installAPIResources(APIGroupPrefix, apiGroupInfo); err != nil {
      return fmt.Errorf("unable to install api resources: %v", err)
    }
  }

  return nil

}

func (s *Server) InstallLegacyAPIGroup(apiPrefix string, apiGroupInfo *APIGroupInfo) error {
  if !s.legacyAPIGroupPrefixes.Has(apiPrefix) {
    return fmt.Errorf("%q not in the allowed legacy API prefixes: %v", apiPrefix, s.legacyAPIGroupPrefixes.List())
  }
  if err := s.installAPIResources(apiPrefix, apiGroupInfo); err != nil {
    return err
  }
  return nil
}

func (s *Server) installAPIResources(apiPrefix string, apiGroupInfo *APIGroupInfo) error {

  for _, groupVersion := range apiGroupInfo.PrioritizedVersions {

    if len(apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version]) == 0 {
      logrus.Warnf("skipping api '%s' because it has no resources", groupVersion)
      continue
    }

    apiGroupVersion := s.getAPIGroupVersion(apiGroupInfo, groupVersion, apiPrefix)

    err := apiGroupVersion.InstallREST(s.Handler.GoRestfulContainer)
    if err != nil {
      return fmt.Errorf("unable to setup api '%v': %v", apiGroupInfo, err)
    }
  }

  return nil
}

func (s *Server) getAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion schema.GroupVersion, apiPrefix string) *endpoints.APIGroupVersion {
  storage := make(map[string]rest.Storage)
  for k, v := range apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version] {
    storage[strings.ToLower(k)] = v
  }
  version := s.newAPIGroupVersion(apiGroupInfo, groupVersion)
  version.Root = apiPrefix
  version.Storage = storage
  return version
}

func (s *Server) newAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion schema.GroupVersion) *endpoints.APIGroupVersion {
  return &endpoints.APIGroupVersion{
    GroupVersion:     groupVersion,
    MetaGroupVersion: apiGroupInfo.MetaGroupVersion,

    ParameterCodec: apiGroupInfo.ParameterCodec,
    Serializer:     apiGroupInfo.NegotiatedSerializer,
    Creater:        apiGroupInfo.Scheme,
    //Convertor:             apiGroupInfo.Scheme,
    //ConvertabilityChecker: apiGroupInfo.Scheme,
    //UnsafeConvertor:       runtime.UnsafeObjectConvertor(apiGroupInfo.Scheme),
    //Defaulter:             apiGroupInfo.Scheme,
    Typer: apiGroupInfo.Scheme,
    //Linker:                runtime.SelfLinker(meta.NewAccessor()),

    //EquivalentResourceRegistry: s.EquivalentResourceRegistry,

    //Admit:             s.admissionControl,
    //MinRequestTimeout: s.minRequestTimeout,
    //Authorizer:        s.Authorizer,
  }
}
