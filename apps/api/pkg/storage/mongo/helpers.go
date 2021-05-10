package mongo

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"reflect"
)

func getNewItemFunc(listObj runtime.Object, v reflect.Value) func() runtime.Object {
	// For unstructured lists with a target group/version, preserve the group/version in the instantiated list items
	//if unstructuredList, isUnstructured := listObj.(*unstructured.UnstructuredList); isUnstructured {
	//  if apiVersion := unstructuredList.GetAPIVersion(); len(apiVersion) > 0 {
	//    return func() runtime.Object {
	//      return &unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": apiVersion}}
	//    }
	//  }
	//}

	// Otherwise just instantiate an empty item
	elem := v.Type().Elem()
	return func() runtime.Object {
		return reflect.New(elem).Interface().(runtime.Object)
	}
}

// appendListItem decodes and appends the object (if it passes filter) to v, which must be a slice.
func appendListItem(v reflect.Value, obj runtime.Object) error {
	//obj, _, err := codec.Decode(data, nil, newItemFunc())
	//if err != nil {
	//  return err
	//}
	// being unable to set the version does not prevent the object from being extracted
	//if err := versioner.UpdateObject(obj, rev); err != nil {
	//  klog.Errorf("failed to update object version: %v", err)
	//}
	//if matched, err := pred.Matches(obj); err == nil && matched {
	//  v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	//}
	v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
	return nil
}
