package runtime

import (
	"bytes"
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkg/conversion/queryparams"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
)

func UseOrCreateObject(t ObjectTyper, c ObjectCreater, gvk schema.GroupVersionKind, obj Object) (Object, error) {
	if obj != nil {
		kinds, _, err := t.ObjectKinds(obj)
		if err != nil {
			return nil, err
		}
		for _, kind := range kinds {
			if gvk == kind {
				return obj, nil
			}
		}
	}
	return c.New(gvk)
}

// Encode is a convenience wrapper for encoding to a []byte from an Encoder
func Encode(e Encoder, obj Object) ([]byte, error) {
	// TODO: reuse buffer
	buf := &bytes.Buffer{}
	if err := e.Encode(obj, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode is a convenience wrapper for decoding data into an Object.
func Decode(d Decoder, data []byte) (Object, error) {
	obj, _, err := d.Decode(data, nil, nil)
	return obj, err
}

// NewParameterCodec creates a ParameterCodec capable of transforming url values into versioned objects and back.
func NewParameterCodec(scheme *Scheme) ParameterCodec {
	return &parameterCodec{
		typer:     scheme,
		convertor: scheme,
		creator:   scheme,
		//defaulter: scheme,
	}
}

// parameterCodec implements conversion to and from query parameters and objects.
type parameterCodec struct {
	typer     ObjectTyper
	convertor ObjectConvertor
	creator   ObjectCreater
	//defaulter ObjectDefaulter
}

var _ ParameterCodec = &parameterCodec{}

// DecodeParameters converts the provided url.Values into an object of type From with the kind of into, and then
// converts that object to into (if necessary). Returns an error if the operation cannot be completed.
func (c *parameterCodec) DecodeParameters(parameters url.Values, from schema.GroupVersion, into Object) error {
	if len(parameters) == 0 {
		return nil
	}
	targetGVKs, _, err := c.typer.ObjectKinds(into)
	if err != nil {
		return err
	}
	for i := range targetGVKs {
		if targetGVKs[i].GroupVersion() == from {
			if err := c.convertor.Convert(&parameters, into, nil); err != nil {
				return err
			}
			// in the case where we going into the same object we're receiving, default on the outbound object
			//if c.defaulter != nil {
			//  c.defaulter.Default(into)
			//}
			return nil
		}
	}

	input, err := c.creator.New(from.WithKind(targetGVKs[0].Kind))
	if err != nil {
		return err
	}
	if err := c.convertor.Convert(&parameters, input, nil); err != nil {
		return err
	}
	// if we have defaulter, default the input before converting to output
	//if c.defaulter != nil {
	//  c.defaulter.Default(input)
	//}
	return c.convertor.Convert(input, into, nil)
}

// EncodeParameters converts the provided object into the to version, then converts that object to url.Values.
// Returns an error if conversion is not possible.
func (c *parameterCodec) EncodeParameters(obj Object, to schema.GroupVersion) (url.Values, error) {
	gvks, _, err := c.typer.ObjectKinds(obj)
	if err != nil {
		return nil, err
	}
	gvk := gvks[0]
	if to != gvk.GroupVersion() {
		out, err := c.convertor.ConvertToVersion(obj, to)
		if err != nil {
			return nil, err
		}
		obj = out
	}
	return queryparams.Convert(obj)
}

var (
	// InternalGroupVersioner will always prefer the internal version for a given group version kind.
	InternalGroupVersioner GroupVersioner = internalGroupVersioner{}
	// DisabledGroupVersioner will reject all kinds passed to it.
	DisabledGroupVersioner GroupVersioner = disabledGroupVersioner{}
)

const (
	internalGroupVersionerIdentifier = "internal"
	disabledGroupVersionerIdentifier = "disabled"
)

type internalGroupVersioner struct{}

// KindForGroupVersionKinds returns an internal Kind if one is found, or converts the first provided kind to the internal version.
func (internalGroupVersioner) KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (schema.GroupVersionKind, bool) {
	for _, kind := range kinds {
		if kind.Version == APIVersionInternal {
			return kind, true
		}
	}
	for _, kind := range kinds {
		return schema.GroupVersionKind{Group: kind.Group, Version: APIVersionInternal, Kind: kind.Kind}, true
	}
	return schema.GroupVersionKind{}, false
}

// Identifier implements GroupVersioner interface.
func (internalGroupVersioner) Identifier() string {
	return internalGroupVersionerIdentifier
}

type disabledGroupVersioner struct{}

// KindForGroupVersionKinds returns false for any input.
func (disabledGroupVersioner) KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (schema.GroupVersionKind, bool) {
	return schema.GroupVersionKind{}, false
}

// Identifier implements GroupVersioner interface.
func (disabledGroupVersioner) Identifier() string {
	return disabledGroupVersionerIdentifier
}

type multiGroupVersioner struct {
	target             schema.GroupVersion
	acceptedGroupKinds []schema.GroupKind
	coerce             bool
}

// NewMultiGroupVersioner returns the provided group version for any kind that matches one of the provided group kinds.
// Kind may be empty in the provided group kind, in which case any kind will match.
func NewMultiGroupVersioner(gv schema.GroupVersion, groupKinds ...schema.GroupKind) GroupVersioner {
	if len(groupKinds) == 0 || (len(groupKinds) == 1 && groupKinds[0].Group == gv.Group) {
		return gv
	}
	return multiGroupVersioner{target: gv, acceptedGroupKinds: groupKinds}
}

// NewCoercingMultiGroupVersioner returns the provided group version for any incoming kind.
// Incoming kinds that match the provided groupKinds are preferred.
// Kind may be empty in the provided group kind, in which case any kind will match.
// Examples:
//   gv=mygroup/__internal, groupKinds=mygroup/Foo, anothergroup/Bar
//   KindForGroupVersionKinds(yetanother/v1/Baz, anothergroup/v1/Bar) -> mygroup/__internal/Bar (matched preferred group/kind)
//
//   gv=mygroup/__internal, groupKinds=mygroup, anothergroup
//   KindForGroupVersionKinds(yetanother/v1/Baz, anothergroup/v1/Bar) -> mygroup/__internal/Bar (matched preferred group)
//
//   gv=mygroup/__internal, groupKinds=mygroup, anothergroup
//   KindForGroupVersionKinds(yetanother/v1/Baz, yetanother/v1/Bar) -> mygroup/__internal/Baz (no preferred group/kind match, uses first kind in list)
func NewCoercingMultiGroupVersioner(gv schema.GroupVersion, groupKinds ...schema.GroupKind) GroupVersioner {
	return multiGroupVersioner{target: gv, acceptedGroupKinds: groupKinds, coerce: true}
}

// KindForGroupVersionKinds returns the target group version if any kind matches any of the original group kinds. It will
// use the originating kind where possible.
func (v multiGroupVersioner) KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (schema.GroupVersionKind, bool) {
	for _, src := range kinds {
		for _, kind := range v.acceptedGroupKinds {
			if kind.Group != src.Group {
				continue
			}
			if len(kind.Kind) > 0 && kind.Kind != src.Kind {
				continue
			}
			return v.target.WithKind(src.Kind), true
		}
	}
	if v.coerce && len(kinds) > 0 {
		return v.target.WithKind(kinds[0].Kind), true
	}
	return schema.GroupVersionKind{}, false
}

// Identifier implements GroupVersioner interface.
func (v multiGroupVersioner) Identifier() string {
	groupKinds := make([]string, 0, len(v.acceptedGroupKinds))
	for _, gk := range v.acceptedGroupKinds {
		groupKinds = append(groupKinds, gk.String())
	}
	result := map[string]string{
		"name":     "multi",
		"target":   v.target.String(),
		"accepted": strings.Join(groupKinds, ","),
		"coerce":   strconv.FormatBool(v.coerce),
	}
	identifier, err := json.Marshal(result)
	if err != nil {
		logrus.Fatalf("Failed marshaling Identifier for %#v: %v", v, err)
	}
	return string(identifier)
}

// SerializerInfoForMediaType returns the first info in types that has a matching media type (which cannot
// include media-type parameters), or the first info with an empty media type, or false if no type matches.
func SerializerInfoForMediaType(types []SerializerInfo, mediaType string) (SerializerInfo, bool) {
	for _, info := range types {
		if info.MediaType == mediaType {
			return info, true
		}
	}
	for _, info := range types {
		if len(info.MediaType) == 0 {
			return info, true
		}
	}
	return SerializerInfo{}, false
}

// codec binds an encoder and decoder.
type codec struct {
	Encoder
	Decoder
}

// NewCodec creates a Codec from an Encoder and Decoder.
func NewCodec(e Encoder, d Decoder) Codec {
	return codec{e, d}
}
