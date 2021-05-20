package handlers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"net/http"
)

type RequestScope struct {
	Serializer runtime.NegotiatedSerializer
	runtime.ParameterCodec
	Creater         runtime.ObjectCreater
	Convertor       runtime.ObjectConvertor
	Defaulter       runtime.ObjectDefaulter
	Typer           runtime.ObjectTyper
	Resource        schema.GroupVersionResource
	Kind            schema.GroupVersionKind
	HubGroupVersion schema.GroupVersion
	Namer           Namer
	UnsafeConvertor runtime.ObjectConvertor
}

func (r *RequestScope) err(err error, w http.ResponseWriter, req *http.Request) {
	responsewriters.ErrorNegotiated(err, r.Serializer, r.Kind.GroupVersion(), w, req)
}

func (r *RequestScope) AllowsMediaTypeTransform(mimeType, mimeSubType string, gvk *schema.GroupVersionKind) bool {
	return mimeType == "application" && mimeSubType == "json"
}

func (r *RequestScope) AllowsServerVersion(version string) bool {
	return version == "v1"
}

func (r *RequestScope) AllowsStreamSchema(s string) bool {
	return false
}

func (r *RequestScope) AcceptsGroupVersion(version schema.GroupVersion) bool {
	return r.Kind.GroupVersion() == version
}
