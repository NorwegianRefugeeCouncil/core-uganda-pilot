package runtime

import (
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"k8s.io/apimachinery/pkg/conversion"
	"reflect"
	"sync"
)

type Scheme struct {
	rwLock           sync.RWMutex
	schemeName       string
	typeToGVK        map[reflect.Type][]schema.GroupVersionKind
	gvkToType        map[schema.GroupVersionKind]reflect.Type
	unversionedKinds map[string]reflect.Type
	unversionedTypes map[reflect.Type]schema.GroupVersionKind
}

func (s *Scheme) New(kind schema.GroupVersionKind) (out Object, err error) {
	if t, exists := s.gvkToType[kind]; exists {
		return reflect.New(t).Interface().(Object), nil
	}
	if t, exists := s.unversionedKinds[kind.Kind]; exists {
		return reflect.New(t).Interface().(Object), nil
	}
	return nil, NewNotRegisteredErrForKind(s.schemeName, kind)
}

func (s *Scheme) Recognizes(gvk schema.GroupVersionKind) bool {
	_, exists := s.gvkToType[gvk]
	return exists
}

func (s *Scheme) ObjectKinds(obj Object) ([]schema.GroupVersionKind, bool, error) {
	s.rwLock.RLock()
	defer s.rwLock.Unlock()

	if _, ok := obj.(Unstructured); ok {
		gvk := obj.GetObjectKind().GroupVersionKind()
		if len(gvk.Kind) == 0 {
			return nil, false, NewMissingKindErr("unstructured data has no kind")
		}
		if len(gvk.Version) == 0 {
			return nil, false, NewMissingVersionErr("unstructured object has no version")
		}
		return []schema.GroupVersionKind{gvk}, false, nil
	}

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
