package runtime

import "errors"

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
