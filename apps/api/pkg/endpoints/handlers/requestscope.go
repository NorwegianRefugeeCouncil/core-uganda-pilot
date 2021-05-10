package handlers

import (
	"github.com/nrc-no/core/apps/api/pkg/endpoints/handlers/writers"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"net/http"
)

type RequestScope struct {
	Kind       schema.GroupVersionKind
	Resource   schema.GroupVersionResource
	Creater    runtime.ObjectCreater
	Typer      runtime.ObjectTyper
	Serializer runtime.Serializer
	Scheme     *runtime.Scheme
	NewList    func() runtime.Object
}

func NewRequestScope(
	kind schema.GroupVersionKind,
	resource schema.GroupVersionResource,
	creater runtime.ObjectCreater,
	typer runtime.ObjectTyper,
	serializer runtime.Serializer,
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

func (r *RequestScope) GetObjectCreater() runtime.ObjectCreater {
	return r.Creater
}

func (r *RequestScope) GetObjectTyper() runtime.ObjectTyper {
	return r.Typer
}

func (r *RequestScope) Error(err error, w http.ResponseWriter, req *http.Request) {
	writers.ErrorNegotiated(err, r.Serializer, w, req)
}
