package endpoints

import (
	"github.com/emicklei/go-restful"
	"github.com/nrc-no/core/apps/api/pkg/registry/rest"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	storageversion2 "github.com/nrc-no/core/apps/api/pkg/storageversion"
	"github.com/nrc-no/core/apps/api/pkg/util/exceptions"
	"k8s.io/apimachinery/pkg/util/sets"
	"path"
	"time"
)

// APIGroupVersion is a helper for exposing rest.Storage objects as http.Handlers via go-restful
// It handles URLs of the form:
// /${storage_key}[/${object_name}]
// Where 'storage_key' points to a rest.Storage object stored in storage.
// This object should contain all parameterization necessary for running a particular API version
type APIGroupVersion struct {
	Storage map[string]rest.Storage

	Root string

	// GroupVersion is the external group version
	GroupVersion schema.GroupVersion

	// OptionsExternalVersion controls the Kubernetes APIVersion used for common objects in the apiserver
	// schema like api.Status, api.DeleteOptions, and metav1.ListOptions. Other implementors may
	// define a version "v1beta1" but want to use the Kubernetes "v1" internal objects. If
	// empty, defaults to GroupVersion.
	OptionsExternalVersion *schema.GroupVersion
	// MetaGroupVersion defaults to "meta.k8s.io/v1" and is the scheme group version used to decode
	// common API implementations like ListOptions. Future changes will allow this to vary by group
	// version (for when the inevitable meta/v2 group emerges).
	MetaGroupVersion *schema.GroupVersion

	// RootScopedKinds are the root scoped kinds for the primary GroupVersion
	RootScopedKinds sets.String

	// Serializer is used to determine how to convert responses from API methods into bytes to send over
	// the wire.
	Serializer     runtime.NegotiatedSerializer
	ParameterCodec runtime.ParameterCodec

	Typer     runtime.ObjectTyper
	Creater   runtime.ObjectCreater
	Convertor runtime.ObjectConvertor
	// ConvertabilityChecker ConvertabilityChecker
	// Defaulter             runtime.ObjectDefaulter
	// Linker                runtime.SelfLinker
	UnsafeConvertor runtime.ObjectConvertor
	// TypeConverter   fieldmanager.TypeConverter

	EquivalentResourceRegistry runtime.EquivalentResourceRegistry

	// Authorizer determines whether a user is allowed to make a certain request. The Handler does a preliminary
	// authorization check using the request URI but it may be necessary to make additional checks, such as in
	// the create-on-update case
	// Authorizer authorizer.Authorizer

	// Admit admission.Interface

	MinRequestTimeout time.Duration

	// OpenAPIModels exposes the OpenAPI models to each individual handler.
	// OpenAPIModels openapiproto.Models

	// The limit on the request body size that would be accepted and decoded in a write request.
	// 0 means no limit.
	MaxRequestBodyBytes int64
}

func (v *APIGroupVersion) InstallREST(container *restful.Container) ([]*storageversion2.ResourceInfo, error) {
	prefix := path.Join(v.Root, v.GroupVersion.Group, v.GroupVersion.Version)
	installer := &APIInstaller{
		group:             v,
		prefix:            prefix,
		minRequestTimeout: v.MinRequestTimeout,
	}

	_, resourceInfos, ws, registrationErrors := installer.Install()

	container.Add(ws)
	return removeNonPersistedResources(resourceInfos), exceptions.NewAggregate(registrationErrors)

}

func removeNonPersistedResources(infos []*storageversion2.ResourceInfo) []*storageversion2.ResourceInfo {
	var filtered []*storageversion2.ResourceInfo
	for _, info := range infos {
		// if EncodingVersion is empty, then the apiserver does not
		// need to register this resource via the storage version API,
		// thus we can remove it.
		if info != nil && len(info.EncodingVersion) > 0 {
			filtered = append(filtered, info)
		}
	}
	return filtered
}
