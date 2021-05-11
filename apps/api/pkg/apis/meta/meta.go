package meta

import (
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
)

// errNotList is returned when an object implements the Object style interfaces but not the List style
// interfaces.
var errNotList = fmt.Errorf("object does not implement the List interfaces")
var errNotCommon = fmt.Errorf("object does not implement the common interface for accessing the SelfLink")

type List metav1.ListInterface
type Type metav1.Type

// ListAccessor returns a List interface for the provided object or an error if the object does
// not provide List.
// IMPORTANT: Objects are NOT a superset of lists. Do not use this check to determine whether an
// object *is* a List.
func ListAccessor(obj interface{}) (List, error) {
	switch t := obj.(type) {
	case List:
		return t, nil
	case ListMetaAccessor:
		if m := t.GetListMeta(); m != nil {
			return m, nil
		}
		return nil, errNotList
	case metav1.ListMetaAccessor:
		if m := t.GetListMeta(); m != nil {
			return m, nil
		}
		return nil, errNotList
	default:
		return nil, errNotList
	}
}
