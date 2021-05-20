package server2

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/api/pkg/server2/registry/rest"
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
	prefix := path.Join(v.Root, v.GroupVersion.Version, v.GroupVersion.Group)
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
