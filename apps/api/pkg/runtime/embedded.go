package runtime

import (
	"errors"
	"github.com/nrc-no/core/apps/api/pkg/conversion"
)

func (e *Unknown) UnmarshalJSON(in []byte) error {
	if e == nil {
		return errors.New("runtime.Unknown: UnmarshalJSON on nil pointer")
	}
	e.TypeMeta = TypeMeta{}
	e.Raw = append(e.Raw[0:0], in...)
	e.ContentEncoding = ""
	e.ContentType = ContentTypeJSON
	return nil
}

// Marshal may get called on pointers or values, so implement MarshalJSON on value.
// http://stackoverflow.com/questions/21390979/custom-marshaljson-never-gets-called-in-go
func (e Unknown) MarshalJSON() ([]byte, error) {
	// If ContentType is unset, we assume this is JSON.
	if e.ContentType != "" && e.ContentType != ContentTypeJSON {
		return nil, errors.New("runtime.Unknown: MarshalJSON on non-json data")
	}
	if e.Raw == nil {
		return []byte("null"), nil
	}
	return e.Raw, nil
}

func Convert_runtime_Object_To_runtime_RawExtension(in *Object, out *RawExtension, s conversion.Scope) error {
	if in == nil {
		out.Raw = []byte("null")
		return nil
	}
	obj := *in
	if unk, ok := obj.(*Unknown); ok {
		if unk.Raw != nil {
			out.Raw = unk.Raw
			return nil
		}
		obj = out.Object
	}
	if obj == nil {
		out.Raw = nil
		return nil
	}
	out.Object = obj
	return nil
}

func Convert_runtime_RawExtension_To_runtime_Object(in *RawExtension, out *Object, s conversion.Scope) error {
	if in.Object != nil {
		*out = in.Object
		return nil
	}
	data := in.Raw
	if len(data) == 0 || (len(data) == 4 && string(data) == "null") {
		*out = nil
		return nil
	}
	*out = &Unknown{
		Raw: data,
		// TODO: Set ContentEncoding and ContentType appropriately.
		// Currently we set ContentTypeJSON to make tests passing.
		ContentType: ContentTypeJSON,
	}
	return nil
}
