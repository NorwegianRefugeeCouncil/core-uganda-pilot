package runtime

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"io"
	"net/url"
)

type Object interface {
	DeepCopyObject() Object
	GetObjectKind() schema.ObjectKind
}

type ResourceVersioner interface {
	SetResourceVersion(obj Object, version int) error
	ResourceVersion(obj Object) (int, error)
}

// GroupVersioner refines a set of possible conversion targets into a single option.
type GroupVersioner interface {
	// KindForGroupVersionKinds returns a desired target group version kind for the given input, or returns ok false if no
	// target is known. In general, if the return target is not in the input list, the caller is expected to invoke
	// Scheme.New(target) and then perform a conversion between the current Go type and the destination Go type.
	// Sophisticated implementations may use additional information about the input kinds to pick a destination kind.
	KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (target schema.GroupVersionKind, ok bool)
	// Identifier returns string representation of the object.
	// Identifiers of two different encoders should be equal only if for every input
	// kinds they return the same result.
	Identifier() string
}

type ObjectCreater interface {
	New(kind schema.GroupVersionKind) (out Object, err error)
}

type ObjectTyper interface {
	ObjectKinds(obj Object) ([]schema.GroupVersionKind, bool, error)
	Recognizes(gvk schema.GroupVersionKind) bool
}

type Encoder interface {
	Encode(obj Object, w io.Writer) error
}

type Decoder interface {
	Decode(data []byte, defaults *schema.GroupVersionKind, into Object) (Object, *schema.GroupVersionKind, error)
}

type Serializer interface {
	Encoder
	Decoder
}

type Codec Serializer

// ParameterCodec defines methods for serializing and deserializing API objects to url.Values and
// performing any necessary conversion. Unlike the normal Codec, query parameters are not self describing
// and the desired version must be specified.
type ParameterCodec interface {
	// DecodeParameters takes the given url.Values in the specified group version and decodes them
	// into the provided object, or returns an error.
	DecodeParameters(parameters url.Values, from schema.GroupVersion, into Object) error
	// EncodeParameters encodes the provided object as query parameters or returns an error.
	EncodeParameters(obj Object, to schema.GroupVersion) (url.Values, error)
}

// Unstructured objects store values as map[string]interface{}, with only values that can be serialized
// to JSON allowed.
type Unstructured interface {
	Object
	// NewEmptyInstance returns a new instance of the concrete type containing only kind/apiVersion and no other data.
	// This should be called instead of reflect.New() for unstructured types because the go type alone does not preserve kind/apiVersion info.
	NewEmptyInstance() Unstructured
	// UnstructuredContent returns a non-nil map with this object's contents. Values may be
	// []interface{}, map[string]interface{}, or any primitive type. Contents are typically serialized to
	// and from JSON. SetUnstructuredContent should be used to mutate the contents.
	UnstructuredContent() map[string]interface{}
	// SetUnstructuredContent updates the object content to match the provided map.
	SetUnstructuredContent(map[string]interface{})
	// IsList returns true if this type is a list or matches the list convention - has an array called "items".
	IsList() bool
	// EachListItem should pass a single item out of the list as an Object to the provided function. Any
	// error should terminate the iteration. If IsList() returns false, this method should return an error
	// instead of calling the provided function.
	EachListItem(func(Object) error) error
}

// ObjectConvertor converts an object to a different version.
type ObjectConvertor interface {
	// Convert attempts to convert one object into another, or returns an error. This
	// method does not mutate the in object, but the in and out object might share data structures,
	// i.e. the out object cannot be mutated without mutating the in object as well.
	// The context argument will be passed to all nested conversions.
	Convert(in, out, context interface{}) error
	// ConvertToVersion takes the provided object and converts it the provided version. This
	// method does not mutate the in object, but the in and out object might share data structures,
	// i.e. the out object cannot be mutated without mutating the in object as well.
	// This method is similar to Convert() but handles specific details of choosing the correct
	// output version.
	ConvertToVersion(in Object, gv GroupVersioner) (out Object, err error)
	ConvertFieldLabel(gvk schema.GroupVersionKind, label, value string) (string, string, error)
}

// NestedObjectDecoder is an optional interface that objects may implement to be given
// an opportunity to decode any nested Objects / RawExtensions during serialization.
type NestedObjectDecoder interface {
	DecodeNestedObjects(d Decoder) error
}

// NestedObjectEncoder is an optional interface that objects may implement to be given
// an opportunity to encode any nested Objects / RawExtensions during serialization.
type NestedObjectEncoder interface {
	EncodeNestedObjects(e Encoder) error
}
