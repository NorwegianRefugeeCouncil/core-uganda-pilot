package json

import (
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"io"
)

type Serializer struct {
	meta    MetaFactory
	creater runtime.ObjectCreater
	typer   runtime.ObjectTyper
}

var _ runtime.Decoder = &Serializer{}
var _ runtime.Encoder = &Serializer{}

func NewSerializer(meta MetaFactory, creater runtime.ObjectCreater, typer runtime.ObjectTyper) *Serializer {
	return &Serializer{meta: meta, creater: creater, typer: typer}
}

func (s *Serializer) Encode(obj runtime.Object, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

func (s *Serializer) Decode(data []byte, gvk *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {

	actualGvk, err := s.meta.Interpret(data)
	if err != nil {
		return nil, nil, err
	}
	if gvk != nil {
		*actualGvk = gvkWithDefaults(*actualGvk, *gvk)
	}

	if into != nil {
		types, _, err := s.typer.ObjectKinds(into)
		switch {
		case err != nil:
			return nil, actualGvk, err
		default:
			*actualGvk = gvkWithDefaults(*actualGvk, types[0])
		}
	}

	if len(actualGvk.Kind) == 0 {
		return nil, actualGvk, runtime.NewMissingKindErr(string(data))
	}
	if len(actualGvk.Version) == 0 {
		return nil, actualGvk, runtime.NewMissingVersionErr(string(data))
	}

	obj, err := runtime.UseOrCreateObject(s.typer, s.creater, *actualGvk, into)
	if err != nil {
		return nil, actualGvk, err
	}

	if err := json.Unmarshal(data, obj); err != nil {
		return nil, actualGvk, err
	}

	return obj, actualGvk, nil

}

func gvkWithDefaults(actual, defaultGVK schema.GroupVersionKind) schema.GroupVersionKind {
	if len(actual.Kind) == 0 {
		actual.Kind = defaultGVK.Kind
	}
	if len(actual.Version) == 0 && len(actual.Group) == 0 {
		actual.Group = defaultGVK.Group
		actual.Version = defaultGVK.Version
	}
	if len(actual.Version) == 0 && actual.Group == defaultGVK.Group {
		actual.Version = defaultGVK.Version
	}
	return actual
}
