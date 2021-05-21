package dynamic

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
)

var watchScheme = runtime.NewScheme()
var basicScheme = runtime.NewScheme()
var deleteScheme = runtime.NewScheme()
var parameterScheme = runtime.NewScheme()
var deleteOptionsCodec = serializer.NewCodecFactory(deleteScheme)
var dynamicParameterCodec = runtime.NewParameterCodec(parameterScheme)

var versionV1 = schema.GroupVersion{Version: "v1"}

func init() {
	metav1.AddToGroupVersion(watchScheme, versionV1)
	metav1.AddToGroupVersion(basicScheme, versionV1)
	metav1.AddToGroupVersion(parameterScheme, versionV1)
	metav1.AddToGroupVersion(deleteScheme, versionV1)
}

// basicNegotiatedSerializer is used to handle discovery and error handling serialization
type basicNegotiatedSerializer struct{}

func (s basicNegotiatedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return []runtime.SerializerInfo{
		{
			MediaType:        "application/json",
			MediaTypeType:    "application",
			MediaTypeSubType: "json",
			EncodesAsText:    true,
			Serializer:       json.NewSerializer(json.DefaultMetaFactory, unstructuredCreater{basicScheme}, unstructuredTyper{basicScheme}, false),
			PrettySerializer: json.NewSerializer(json.DefaultMetaFactory, unstructuredCreater{basicScheme}, unstructuredTyper{basicScheme}, true),
			StreamSerializer: &runtime.StreamSerializerInfo{
				EncodesAsText: true,
				Serializer:    json.NewSerializer(json.DefaultMetaFactory, basicScheme, basicScheme, false),
				Framer:        json.Framer,
			},
		},
	}
}

func (s basicNegotiatedSerializer) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return runtime.WithVersionEncoder{
		Version:     gv,
		Encoder:     encoder,
		ObjectTyper: unstructuredTyper{basicScheme},
	}
}

func (s basicNegotiatedSerializer) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return decoder
}

type unstructuredCreater struct {
	nested runtime.ObjectCreater
}

func (c unstructuredCreater) New(kind schema.GroupVersionKind) (runtime.Object, error) {
	out, err := c.nested.New(kind)
	if err == nil {
		return out, nil
	}
	out = &unstructured.Unstructured{}
	out.GetObjectKind().SetGroupVersionKind(kind)
	return out, nil
}

type unstructuredTyper struct {
	nested runtime.ObjectTyper
}

func (t unstructuredTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	kinds, unversioned, err := t.nested.ObjectKinds(obj)
	if err == nil {
		return kinds, unversioned, nil
	}
	if _, ok := obj.(runtime.Unstructured); ok && !obj.GetObjectKind().GroupVersionKind().Empty() {
		return []schema.GroupVersionKind{obj.GetObjectKind().GroupVersionKind()}, false, nil
	}
	return nil, false, err
}

func (t unstructuredTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}
