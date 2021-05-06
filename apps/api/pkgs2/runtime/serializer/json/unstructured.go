package json

import (
	"encoding/json"
	"github.com/nrc-no/core/apps/api/pkgs2/apis/meta/v1/unstructured"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"io"
	"strings"
)

// UnstructuredJSONScheme is capable of converting JSON data into the Unstructured
// type, which can be used for generic access to objects without a predefined scheme.
// TODO: move into serializer/json.
// var UnstructuredJSONScheme runtime.Codec = unstructuredJSONScheme{}

var UnstructuredJSONScheme = unstructuredJSONScheme{}

type unstructuredJSONScheme struct{}

const unstructuredJSONSchemeIdentifier runtime.Identifier = "unstructuredJSON"

func (s unstructuredJSONScheme) Decode(data []byte, _ *schema.GroupVersionKind, obj runtime.Object) (runtime.Object, *schema.GroupVersionKind, error) {
	var err error
	if obj != nil {
		err = s.decodeInto(data, obj)
	} else {
		obj, err = s.decode(data)
	}

	if err != nil {
		return nil, nil, err
	}

	gvk := obj.GetObjectKind().GroupVersionKind()
	if len(gvk.Kind) == 0 {
		return nil, &gvk, runtime.NewMissingKindErr(string(data))
	}

	return obj, &gvk, nil
}

func (s unstructuredJSONScheme) Encode(obj runtime.Object, w io.Writer) error {
	//if co, ok := obj.(runtime.CacheableObject); ok {
	//  return co.CacheEncode(s.Identifier(), s.doEncode, w)
	//}
	return s.doEncode(obj, w)
}

func (unstructuredJSONScheme) doEncode(obj runtime.Object, w io.Writer) error {
	switch t := obj.(type) {
	case *unstructured.Unstructured:
		bytes, err := json.Marshal(t.Object)
		if err != nil {
			return err
		}
		if _, err := w.Write(bytes); err != nil {
			return err
		}
		return nil
	case *unstructured.UnstructuredList:
		items := make([]interface{}, 0, len(t.Items))
		for _, i := range t.Items {
			items = append(items, i.Object)
		}
		listObj := make(map[string]interface{}, len(t.Object)+1)
		for k, v := range t.Object { // Make a shallow copy
			listObj[k] = v
		}
		listObj["items"] = items
		bytes, err := json.Marshal(listObj)
		if err != nil {
			return err
		}
		if _, err := w.Write(bytes); err != nil {
			return err
		}
		return nil
	//case *runtime.Unknown:
	//  // TODO: Unstructured needs to deal with ContentType.
	//  _, err := w.Write(t.Raw)
	//  return err
	default:
		bytes, err := json.Marshal(t)
		if err != nil {
			return err
		}
		if _, err := w.Write(bytes); err != nil {
			return err
		}
		return nil
	}
}

// Identifier implements runtime.Encoder interface.
func (unstructuredJSONScheme) Identifier() runtime.Identifier {
	return unstructuredJSONSchemeIdentifier
}

func (s unstructuredJSONScheme) decode(data []byte) (runtime.Object, error) {
	type detector struct {
		Items json.RawMessage
	}
	var det detector
	if err := json.Unmarshal(data, &det); err != nil {
		return nil, err
	}

	if det.Items != nil {
		list := &unstructured.UnstructuredList{}
		err := s.decodeToList(data, list)
		return list, err
	}

	// No Items field, so it wasn't a list.
	unstruct := &unstructured.Unstructured{}
	err := s.decodeToUnstructured(data, unstruct)
	return unstruct, err
}

func (s unstructuredJSONScheme) decodeInto(data []byte, obj runtime.Object) error {
	switch x := obj.(type) {
	case *unstructured.Unstructured:
		return s.decodeToUnstructured(data, x)
	case *unstructured.UnstructuredList:
		return s.decodeToList(data, x)
	default:
		return json.Unmarshal(data, x)
	}
}

func (unstructuredJSONScheme) decodeToUnstructured(data []byte, unstruct *unstructured.Unstructured) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	unstruct.Object = m

	return nil
}

func (s unstructuredJSONScheme) decodeToList(data []byte, list *unstructured.UnstructuredList) error {
	type decodeList struct {
		Items []json.RawMessage
	}

	var dList decodeList
	if err := json.Unmarshal(data, &dList); err != nil {
		return err
	}

	if err := json.Unmarshal(data, &list.Object); err != nil {
		return err
	}

	// For typed lists, e.g., a PodList, API server doesn't set each item's
	// APIVersion and Kind. We need to set it.
	listAPIVersion := list.GetAPIVersion()
	listKind := list.GetKind()
	itemKind := strings.TrimSuffix(listKind, "List")

	delete(list.Object, "items")
	list.Items = make([]unstructured.Unstructured, 0, len(dList.Items))
	for _, i := range dList.Items {
		unstruct := &unstructured.Unstructured{}
		if err := s.decodeToUnstructured([]byte(i), unstruct); err != nil {
			return err
		}
		// This is hacky. Set the item's Kind and APIVersion to those inferred
		// from the List.
		if len(unstruct.GetKind()) == 0 && len(unstruct.GetAPIVersion()) == 0 {
			unstruct.SetKind(itemKind)
			unstruct.SetAPIVersion(listAPIVersion)
		}
		list.Items = append(list.Items, *unstruct)
	}
	return nil
}
