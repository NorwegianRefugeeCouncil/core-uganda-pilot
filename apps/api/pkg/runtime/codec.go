package runtime

import (
	"bytes"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
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
