package json

import (
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
)

// SerializerOptions holds the options which are used to configure a JSON/YAML serializer.
// example:
// (1) To configure a JSON serializer, set `Yaml` to `false`.
// (2) To configure a YAML serializer, set `Yaml` to `true`.
// (3) To configure a strict serializer that can return strictDecodingError, set `Strict` to `true`.
type SerializerOptions struct {
	// Yaml: configures the Serializer to work with JSON(false) or YAML(true).
	// When `Yaml` is enabled, this serializer only supports the subset of YAML that
	// matches JSON, and will error if constructs are used that do not serialize to JSON.
	Yaml bool

	// Pretty: configures a JSON enabled Serializer(`Yaml: false`) to produce human-readable output.
	// This option is silently ignored when `Yaml` is `true`.
	Pretty bool

	// Strict: configures the Serializer to return strictDecodingError's when duplicate fields are present decoding JSON or YAML.
	// Note that enabling this option is not as performant as the non-strict variant, and should not be used in fast paths.
	Strict bool
}

type Serializer struct {
	meta       MetaFactory
	creater    runtime.ObjectCreater
	typer      runtime.ObjectTyper
	options    SerializerOptions
	identifier runtime.Identifier
}

var _ runtime.Decoder = &Serializer{}
var _ runtime.Encoder = &Serializer{}

func NewSerializer(meta MetaFactory, creater runtime.ObjectCreater, typer runtime.ObjectTyper, pretty bool) *Serializer {
	return NewSerializerWithOptions(meta, creater, typer, SerializerOptions{false, pretty, false})
}

func NewYAMLSerializer(meta MetaFactory, creater runtime.ObjectCreater, typer runtime.ObjectTyper) *Serializer {
	return NewSerializerWithOptions(meta, creater, typer, SerializerOptions{true, false, false})
}

func NewSerializerWithOptions(meta MetaFactory, creater runtime.ObjectCreater, typer runtime.ObjectTyper, options SerializerOptions) *Serializer {
	return &Serializer{
		meta:       meta,
		creater:    creater,
		typer:      typer,
		options:    options,
		identifier: identifier(options),
	}
}

// identifier computes Identifier of Encoder based on the given options.
func identifier(options SerializerOptions) runtime.Identifier {
	result := map[string]string{
		"name":   "json",
		"yaml":   strconv.FormatBool(options.Yaml),
		"pretty": strconv.FormatBool(options.Pretty),
	}
	identifier, err := json.Marshal(result)
	if err != nil {
		logrus.Fatalf("Failed marshaling identifier for json Serializer: %v", err)
	}
	return runtime.Identifier(identifier)
}

func (s *Serializer) Encode(obj runtime.Object, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(obj)
}

func (s *Serializer) Decode(originalData []byte, gvk *schema.GroupVersionKind, into runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {

	data := originalData

	actualGvk, err := s.meta.Interpret(data)
	if err != nil {
		return nil, nil, err
	}
	if gvk != nil {
		*actualGvk = gvkWithDefaults(*actualGvk, *gvk)
	}

	if unk, ok := into.(*runtime.Unknown); ok && unk != nil {
		unk.Raw = originalData
		unk.ContentType = runtime.ContentTypeJSON
		unk.GetObjectKind().SetGroupVersionKind(*actualGvk)
		return unk, actualGvk, nil
	}

	if into != nil {

		_, isUnstructured := into.(runtime.Unstructured)
		types, _, err := s.typer.ObjectKinds(into)
		switch {
		case runtime.IsNotRegisteredError(err), isUnstructured:
			if err := json.Unmarshal(data, into); err != nil {
				return nil, actualGvk, err
			}
			return into, actualGvk, nil
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
