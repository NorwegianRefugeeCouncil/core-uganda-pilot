package handlers

import (
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/writers"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RequestScope struct {
	Kind                     schema.GroupVersionKind
	Resource                 schema.GroupVersionResource
	Creater                  runtime.ObjectCreater
	Typer                    runtime.ObjectTyper
	Serializer               runtime.NegotiatedSerializer
	StandardSerializers      []runtime.SerializerInfo
	Scheme                   *runtime.Scheme
	NewList                  func() runtime.Object
	ParameterCodec           runtime.ParameterCodec
	Convertor                runtime.ObjectConvertor
	EquivalentResourceMapper runtime.EquivalentResourceRegistry
	HubGroupVersion          schema.GroupVersion
	Namer                    ScopeNamer
}

func NewRequestScope(
	kind schema.GroupVersionKind,
	resource schema.GroupVersionResource,
	creater runtime.ObjectCreater,
	typer runtime.ObjectTyper,
	serializer runtime.NegotiatedSerializer,
	Scheme *runtime.Scheme,
	newList func() runtime.Object,
) *RequestScope {
	return &RequestScope{
		Kind:       kind,
		Resource:   resource,
		Creater:    creater,
		Typer:      typer,
		Serializer: serializer,
		Scheme:     Scheme,
		NewList:    newList,
	}
}

func (r *RequestScope) Error(err error, w http.ResponseWriter, req *http.Request) {
	logrus.Errorf("request error: %v", err)
	writers.ErrorNegotiated(err, r.Serializer, r.Kind.GroupVersion(), w, req)
}

// AcceptsGroupVersion returns true if the specified GroupVersion is allowed
// in create and update requests.
func (scope *RequestScope) AcceptsGroupVersion(gv schema.GroupVersion) bool {
	// If there's a custom acceptor, delegate to it. This is extremely rare.
	// if scope.AcceptsGroupVersionDelegate != nil {
	//  return scope.AcceptsGroupVersionDelegate.AcceptsGroupVersion(gv)
	//}
	// Fall back to only allowing the singular Kind. This is the typical behavior.
	return gv == scope.Kind.GroupVersion()
}

func (scope *RequestScope) AllowsMediaTypeTransform(mimeType, mimeSubType string, gvk *schema.GroupVersionKind) bool {
	// some handlers like CRDs can't serve all the mime types that PartialObjectMetadata or Table can - if
	// gvk is nil (no conversion) allow StandardSerializers to further restrict the set of mime types.
	if gvk == nil {
		if len(scope.StandardSerializers) == 0 {
			return true
		}
		for _, info := range scope.StandardSerializers {
			if info.MediaTypeType == mimeType && info.MediaTypeSubType == mimeSubType {
				return true
			}
		}
		return false
	}

	// TODO: this is temporary, replace with an abstraction calculated at endpoint installation time
	if gvk.GroupVersion() == metav1.SchemeGroupVersion {
		switch gvk.Kind {
		//case "Table":
		//  return scope.TableConvertor != nil &&
		//    mimeType == "application" &&
		//    (mimeSubType == "json" || mimeSubType == "yaml")
		//case "PartialObjectMetadata", "PartialObjectMetadataList":
		//  // TODO: should delineate between lists and non-list endpoints
		//  return true
		default:
			return false
		}
	}
	return false
}

func (scope *RequestScope) AllowsServerVersion(version string) bool {
	return true
	// return version == scope.MetaGroupVersion.Version
}

func (scope *RequestScope) AllowsStreamSchema(s string) bool {
	return s == "watch"
}

// var _ admission.ObjectInterfaces = &RequestScope{}

func (r *RequestScope) GetObjectCreater() runtime.ObjectCreater { return r.Creater }
func (r *RequestScope) GetObjectTyper() runtime.ObjectTyper     { return r.Typer }

// func (r *RequestScope) GetObjectDefaulter() runtime.ObjectDefaulter { return r.Defaulter }
func (r *RequestScope) GetObjectConvertor() runtime.ObjectConvertor { return r.Convertor }
func (r *RequestScope) GetEquivalentResourceMapper() runtime.EquivalentResourceMapper {
	return r.EquivalentResourceMapper
}
