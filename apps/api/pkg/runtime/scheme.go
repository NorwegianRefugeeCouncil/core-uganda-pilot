package runtime

import (
	"fmt"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
	"github.com/nrc-no/core/apps/api/pkg/util/conversion"
	"reflect"
)

const (
	// APIVersionInternal may be used if you are registering a type that should not
	// be considered stable or serialized - it is a convention only and has no
	// special behavior in this package.
	APIVersionInternal = "__internal"
)

type Scheme struct {
	schemeName       string
	gvkToType        map[schema.GroupVersionKind]reflect.Type
	typeToGVK        map[reflect.Type][]schema.GroupVersionKind
	unversionedTypes map[reflect.Type]schema.GroupVersionKind
	observedVersions []schema.GroupVersion
}

var _ ObjectTyper = &Scheme{}
var _ ObjectCreater = &Scheme{}

func NewScheme() *Scheme {
	return &Scheme{
		gvkToType:        map[schema.GroupVersionKind]reflect.Type{},
		typeToGVK:        map[reflect.Type][]schema.GroupVersionKind{},
		unversionedTypes: map[reflect.Type]schema.GroupVersionKind{},
		observedVersions: []schema.GroupVersion{},
	}
}

func (s *Scheme) New(kind schema.GroupVersionKind) (Object, error) {
	if t, exists := s.gvkToType[kind]; exists {
		return reflect.New(t).Interface().(Object), nil
	}

	//if t, exists := s.unversionedKinds[kind.Kind]; exists {
	//  return reflect.New(t).Interface().(Object), nil
	//}
	return nil, NewNotRegisteredErrForKind(s.schemeName, kind)
}

// ObjectKinds returns all possible group,version,kind of the go object, true if the
// object is considered unversioned, or an error if it's not a pointer or is unregistered.
func (s *Scheme) ObjectKinds(obj Object) ([]schema.GroupVersionKind, bool, error) {
	// Unstructured objects are always considered to have their declared GVK
	//if _, ok := obj.(Unstructured); ok {
	//  // we require that the GVK be populated in order to recognize the object
	//  gvk := obj.GetObjectKind().GroupVersionKind()
	//  if len(gvk.Kind) == 0 {
	//    return nil, false, NewMissingKindErr("unstructured object has no kind")
	//  }
	//  if len(gvk.Version) == 0 {
	//    return nil, false, NewMissingVersionErr("unstructured object has no version")
	//  }
	//  return []schema.GroupVersionKind{gvk}, false, nil
	//}

	v, err := conversion.EnforcePtr(obj)
	if err != nil {
		return nil, false, err
	}
	t := v.Type()

	gvks, ok := s.typeToGVK[t]
	if !ok {
		return nil, false, NewNotRegisteredErrForType(s.schemeName, t)
	}
	_, unversionedType := s.unversionedTypes[t]

	return gvks, unversionedType, nil
}

// Recognizes returns true if the scheme is able to handle the provided group,version,kind
// of an object.
func (s *Scheme) Recognizes(gvk schema.GroupVersionKind) bool {
	_, exists := s.gvkToType[gvk]
	return exists
}

// AddKnownTypes registers all types passed in 'types' as being members of version 'version'.
// All objects passed to types should be pointers to structs. The name that go reports for
// the struct becomes the "kind" field when encoding. Version may not be empty - use the
// APIVersionInternal constant if you have a type that does not have a formal version.
func (s *Scheme) AddKnownTypes(gv schema.GroupVersion, types ...Object) {
	s.addObservedVersion(gv)
	for _, obj := range types {
		t := reflect.TypeOf(obj)
		if t.Kind() != reflect.Ptr {
			panic("All types must be pointers to structs.")
		}
		t = t.Elem()
		s.AddKnownTypeWithName(gv.WithKind(t.Name()), obj)
	}
}

// AddKnownTypeWithName is like AddKnownTypes, but it lets you specify what this type should
// be encoded as. Useful for testing when you don't want to make multiple packages to define
// your structs. Version may not be empty - use the APIVersionInternal constant if you have a
// type that does not have a formal version.
func (s *Scheme) AddKnownTypeWithName(gvk schema.GroupVersionKind, obj Object) {
	s.addObservedVersion(gvk.GroupVersion())
	t := reflect.TypeOf(obj)
	if len(gvk.Version) == 0 {
		panic(fmt.Sprintf("version is required on all types: %s %v", gvk, t))
	}
	if t.Kind() != reflect.Ptr {
		panic("All types must be pointers to structs.")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("All types must be pointers to structs.")
	}

	if oldT, found := s.gvkToType[gvk]; found && oldT != t {
		panic(fmt.Sprintf("Double registration of different types for %v: old=%v.%v, new=%v.%v in scheme %q", gvk, oldT.PkgPath(), oldT.Name(), t.PkgPath(), t.Name(), s.schemeName))
	}

	s.gvkToType[gvk] = t

	for _, existingGvk := range s.typeToGVK[t] {
		if existingGvk == gvk {
			return
		}
	}
	s.typeToGVK[t] = append(s.typeToGVK[t], gvk)

	// if the type implements DeepCopyInto(<obj>), register a self-conversion
	//if m := reflect.ValueOf(obj).MethodByName("DeepCopyInto"); m.IsValid() && m.Type().NumIn() == 1 && m.Type().NumOut() == 0 && m.Type().In(0) == reflect.TypeOf(obj) {
	//  if err := s.AddGeneratedConversionFunc(obj, obj, func(a, b interface{}, scope conversion.Scope) error {
	//    // copy a to b
	//    reflect.ValueOf(a).MethodByName("DeepCopyInto").Call([]reflect.Value{reflect.ValueOf(b)})
	//    // clear TypeMeta to match legacy reflective conversion
	//    b.(Object).GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{})
	//    return nil
	//  }); err != nil {
	//    panic(err)
	//  }
	//}
}

func (s *Scheme) addObservedVersion(version schema.GroupVersion) {
	if len(version.Version) == 0 || version.Version == APIVersionInternal {
		return
	}
	for _, observedVersion := range s.observedVersions {
		if observedVersion == version {
			return
		}
	}

	s.observedVersions = append(s.observedVersions, version)
}
