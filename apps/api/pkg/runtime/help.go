package runtime

import (
	"errors"
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/conversion"
	"io"
	"reflect"
)

var (
	errExpectFieldItems = errors.New("no Items field in this object")
	errExpectSliceItems = errors.New("Items field must be a slice of objects")
)

func GetItemsPtr(list Object) (interface{}, error) {
	obj, err := getItemsPtr(list)
	if err != nil {
		return nil, fmt.Errorf("%T is not a list: %v", list, err)
	}
	return obj, nil
}

// getItemsPtr returns a pointer to the list object's Items member or an error.
func getItemsPtr(list Object) (interface{}, error) {
	v, err := conversion.EnforcePtr(list)
	if err != nil {
		return nil, err
	}

	items := v.FieldByName("Items")
	if !items.IsValid() {
		return nil, errExpectFieldItems
	}
	switch items.Kind() {
	case reflect.Interface, reflect.Ptr:
		target := reflect.TypeOf(items.Interface()).Elem()
		if target.Kind() != reflect.Slice {
			return nil, errExpectSliceItems
		}
		return items.Interface(), nil
	case reflect.Slice:
		return items.Addr().Interface(), nil
	default:
		return nil, errExpectSliceItems
	}
}

// unsafeObjectConvertor implements ObjectConvertor using the unsafe conversion path.
type unsafeObjectConvertor struct {
	*Scheme
}

var _ ObjectConvertor = unsafeObjectConvertor{}

// ConvertToVersion converts in to the provided outVersion without copying the input first, which
// is only safe if the output object is not mutated or reused.
func (c unsafeObjectConvertor) ConvertToVersion(in Object, outVersion GroupVersioner) (Object, error) {
	return c.Scheme.UnsafeConvertToVersion(in, outVersion)
}

// UnsafeObjectConvertor performs object conversion without copying the object structure,
// for use when the converted object will not be reused or mutated. Primarily for use within
// versioned codecs, which use the external object for serialization but do not return it.
func UnsafeObjectConvertor(scheme *Scheme) ObjectConvertor {
	return unsafeObjectConvertor{scheme}
}

// WithVersionEncoder serializes an object and ensures the GVK is set.
type WithVersionEncoder struct {
	Version GroupVersioner
	Encoder
	ObjectTyper
}

// Encode does not do conversion. It sets the gvk during serialization.
func (e WithVersionEncoder) Encode(obj Object, stream io.Writer) error {
	gvks, _, err := e.ObjectTyper.ObjectKinds(obj)
	if err != nil {
		if IsNotRegisteredError(err) {
			return e.Encoder.Encode(obj, stream)
		}
		return err
	}
	kind := obj.GetObjectKind()
	oldGVK := kind.GroupVersionKind()
	gvk := gvks[0]
	if e.Version != nil {
		preferredGVK, ok := e.Version.KindForGroupVersionKinds(gvks)
		if ok {
			gvk = preferredGVK
		}
	}
	kind.SetGroupVersionKind(gvk)
	err = e.Encoder.Encode(obj, stream)
	kind.SetGroupVersionKind(oldGVK)
	return err
}

// WithoutVersionDecoder clears the group version kind of a deserialized object.
type WithoutVersionDecoder struct {
	Decoder
}

// Decode does not do conversion. It removes the gvk during deserialization.
func (d WithoutVersionDecoder) Decode(data []byte, defaults *schema.GroupVersionKind, into Object) (Object, *schema.GroupVersionKind, error) {
	obj, gvk, err := d.Decoder.Decode(data, defaults, into)
	if obj != nil {
		kind := obj.GetObjectKind()
		// clearing the gvk is just a convention of a codec
		kind.SetGroupVersionKind(schema.GroupVersionKind{})
	}
	return obj, gvk, err
}

func SetZeroValue(objPtr Object) error {
	v, err := conversion.EnforcePtr(objPtr)
	if err != nil {
		return err
	}
	v.Set(reflect.Zero(v.Type()))
	return nil
}
