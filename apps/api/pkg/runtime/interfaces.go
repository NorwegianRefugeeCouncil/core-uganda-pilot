package runtime

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"io"
)

type Object interface {
	GetResourceVersion() int
	SetResourceVersion(version int)
	GetAPIVersion() string
	SetAPIVersion(version string)
	GetAPIGroup() string
	SetAPIGroup(apiGroup string)
	GetKind() string
	SetKind(kind string)
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
