package server2

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"path"
	"strings"
)

type APIGroupInfo struct {
	PrioritizedVersions          []schema.GroupVersion
	VersionedResourcesStorageMap map[string]map[string]rest.Storage
	Scheme                       *runtime.Scheme
	Serializer                   runtime.NegotiatedSerializer
	ParameterCodec               runtime.ParameterCodec
}

func NewDefaultAPIGroup(group string, scheme *runtime.Scheme, parameterCodec runtime.ParameterCodec, codecs serializer.CodecFactory) APIGroupInfo {
	return APIGroupInfo{
		PrioritizedVersions:          scheme.PrioritizedVersionsForGroup(group),
		VersionedResourcesStorageMap: map[string]map[string]rest.Storage{},
		Scheme:                       scheme,
		ParameterCodec:               parameterCodec,
		Serializer:                   codecs,
	}
}

func (a *APIGroupInfo) GetAPIGroupVersion(groupVersion schema.GroupVersion, apiPrefix string) *APIGroupVersion {
	storage := make(map[string]rest.Storage)
	for k, v := range a.VersionedResourcesStorageMap[groupVersion.Version] {
		storage[strings.ToLower(k)] = v
	}
	version := NewAPIGroupVersion(a, groupVersion)
	version.Storage = storage
	version.Root = apiPrefix
	return version
}

type APIGroupVersion struct {
	Storage        map[string]rest.Storage
	GroupVersion   schema.GroupVersion
	Serializer     runtime.NegotiatedSerializer
	ParameterCodec runtime.ParameterCodec
	Typer          runtime.ObjectTyper
	Creater        runtime.ObjectCreater
	Convertor      runtime.ObjectConvertor
	Defaulter      runtime.ObjectDefaulter
	Root           string
}

func (v *APIGroupVersion) InstallREST(container *restful.Container) error {
	prefix := path.Join(v.Root, v.GroupVersion.Group, v.GroupVersion.Version)
	installer := &APIInstaller{
		group:  v,
		prefix: prefix,
	}
	ws, err := installer.Install()
	if err != nil {
		return err
	}
	container.Add(ws)
	return nil
}

func NewAPIGroupVersion(apiGroupInfo *APIGroupInfo, groupVersion schema.GroupVersion) *APIGroupVersion {
	return &APIGroupVersion{
		GroupVersion:   groupVersion,
		Serializer:     apiGroupInfo.Serializer,
		ParameterCodec: apiGroupInfo.ParameterCodec,
		Typer:          apiGroupInfo.Scheme,
		Creater:        apiGroupInfo.Scheme,
		Convertor:      apiGroupInfo.Scheme,
		Defaulter:      apiGroupInfo.Scheme,
	}
}

// installApiGroups registers the API groups into go-restful container
// this method will register the necessary routes and handlers
func installApiGroups(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfos ...*APIGroupInfo) error {

	for _, apiGroupInfo := range apiGroupInfos {
		if len(apiGroupInfo.PrioritizedVersions[0].Group) == 0 {
			return fmt.Errorf("cannot register handler with an empty group for %#v", *apiGroupInfo)
		}
		if len(apiGroupInfo.PrioritizedVersions[0].Version) == 0 {
			return fmt.Errorf("cannot register handler with an empty version for %#v", *apiGroupInfo)
		}
	}

	for _, apiGroupInfo := range apiGroupInfos {
		if err := installApiResources(goRestfulContainer, apiPrefix, apiGroupInfo); err != nil {
			return err
		}
	}

	return nil
}

func installApiResources(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfo *APIGroupInfo) error {
	for _, groupVersion := range apiGroupInfo.PrioritizedVersions {

		if len(apiGroupInfo.VersionedResourcesStorageMap[groupVersion.Version]) == 0 {
			logrus.Warnf("skipping api %v because it has no resources", groupVersion)
			continue
		}

		apiGroupVersion := apiGroupInfo.GetAPIGroupVersion(groupVersion, apiPrefix)
		if err := apiGroupVersion.InstallREST(goRestfulContainer); err != nil {
			return err
		}

	}
	return nil
}

// installApiGroup registers an API group into the go-restful container.
// see installApiGroups
func installApiGroup(goRestfulContainer *restful.Container, apiPrefix string, apiGroupInfo *APIGroupInfo) error {
	return installApiGroups(goRestfulContainer, apiPrefix, apiGroupInfo)
}
