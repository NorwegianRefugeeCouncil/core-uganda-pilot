package endpoints

import (
	"github.com/emicklei/go-restful"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ConvertabilityChecker interface {
	VersionsForGroupKind(gk schema.GroupKind) []schema.GroupVersion
}

type FormDefinitionHandler struct {
}

func Install(
	resourceName string,

	container *restful.Container,
	groupVersion schema.GroupVersion,
	serializer runtime.NegotiatedSerializer,
	parameterCodec runtime.ParameterCodec,
	typer runtime.ObjectTyper,
	creater runtime.ObjectCreater,
	convertor runtime.ObjectConvertor,
	convertabilityChecker ConvertabilityChecker,
	defaulter runtime.ObjectDefaulter,
) {

}
