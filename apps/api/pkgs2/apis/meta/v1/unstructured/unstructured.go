package unstructured

import (
	"errors"
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkgs2/apis/meta/v1"
	runtime2 "github.com/nrc-no/core/apps/api/pkgs2/runtime"
	"github.com/nrc-no/core/apps/api/pkgs2/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type Unstructured struct {
	Object map[string]interface{}
}

var _ metav1.Object = &Unstructured{}
var _ runtime2.Unstructured = &Unstructured{}
var _ metav1.ListInterface = &Unstructured{}

func (u *Unstructured) GetName() string {
	return getNestedString(u.Object, "metadata", "name")
}

func (u *Unstructured) SetName(name string) {
	if len(name) == 0 {
		RemoveNestedField(u.Object, "metadata", "name")
		return
	}
	u.setNestedField(name, "metadata", "name")
}

func (u *Unstructured) GetUID() string {
	return getNestedString(u.Object, "metadata", "uid")
}

func (u *Unstructured) SetUID(uid string) {
	if len(string(uid)) == 0 {
		RemoveNestedField(u.Object, "metadata", "uid")
		return
	}
	u.setNestedField(string(uid), "metadata", "uid")
}

func (u *Unstructured) GetResourceVersion() string {
	return getNestedString(u.Object, "metadata", "resourceVersion")
}

func (u *Unstructured) SetResourceVersion(version string) {
	if len(version) == 0 {
		RemoveNestedField(u.Object, "metadata", "resourceVersion")
		return
	}
	u.setNestedField(version, "metadata", "resourceVersion")
}

func (u *Unstructured) GetCreationTimestamp() time.Time {
	timestamp, _ := time.Parse(time.RFC3339, getNestedString(u.Object, "metadata", "creationTimestamp"))
	return timestamp
}

func (u *Unstructured) SetCreationTimestamp(timestamp time.Time) {
	if timestamp.IsZero() {
		RemoveNestedField(u.Object, "metadata", "creationTimestamp")
		return
	}
	u.setNestedField(timestamp, "metadata", "creationTimestamp")
}

func (u *Unstructured) GetDeletionTimestamp() *time.Time {
	timestamp, _ := time.Parse(time.RFC3339, getNestedString(u.Object, "metadata", "deletionTimestamp"))
	if timestamp.IsZero() {
		return nil
	}
	return &timestamp
}

func (u *Unstructured) SetDeletionTimestamp(timestamp *time.Time) {
	if timestamp == nil {
		RemoveNestedField(u.Object, "metadata", "deletionTimestamp")
		return
	}
	u.setNestedField(timestamp, "metadata", "creationTimestamp")
}

func (u *Unstructured) GetLabels() map[string]string {
	m, _, _ := NestedStringMap(u.Object, "metadata", "labels")
	return m
}

func (u *Unstructured) SetLabels(labels map[string]string) {
	if labels == nil {
		RemoveNestedField(u.Object, "metadata", "labels")
		return
	}
	u.setNestedMap(labels, "metadata", "labels")
}

func (u *Unstructured) GetAnnotations() map[string]string {
	m, _, _ := NestedStringMap(u.Object, "metadata", "annotations")
	return m
}

func (u *Unstructured) SetAnnotations(annotations map[string]string) {
	if annotations == nil {
		RemoveNestedField(u.Object, "metadata", "annotations")
		return
	}
	u.setNestedMap(annotations, "metadata", "annotations")
}

func (u *Unstructured) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	u.SetAPIVersion(gvk.GroupVersion().String())
	u.SetKind(gvk.Kind)
}

func (u *Unstructured) GroupVersionKind() schema.GroupVersionKind {
	gv, err := schema.ParseGroupVersion(u.GetAPIVersion())
	if err != nil {
		return schema.GroupVersionKind{}
	}
	gvk := gv.WithKind(u.GetKind())
	return gvk
}

func (u *Unstructured) GetObjectKind() schema.ObjectKind {
	return u
}

func (u *Unstructured) DeepCopyObject() runtime2.Object {
	if c := u.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (u *Unstructured) DeepCopyInto(out *Unstructured) {
	clone := u.DeepCopy()
	*out = *clone
	return
}

func (u *Unstructured) NewEmptyInstance() runtime2.Unstructured {
	out := new(Unstructured)
	if u != nil {
		out.GetObjectKind().SetGroupVersionKind(u.GetObjectKind().GroupVersionKind())
	}
	return out
}

func (u *Unstructured) UnstructuredContent() map[string]interface{} {
	if u.Object == nil {
		return make(map[string]interface{})
	}
	return u.Object
}

func (u *Unstructured) SetUnstructuredContent(m map[string]interface{}) {
	u.Object = m
}

func (u *Unstructured) IsList() bool {
	field, ok := u.Object["items"]
	if !ok {
		return false
	}
	_, ok = field.([]interface{})
	return ok
}

func (u *Unstructured) EachListItem(fn func(runtime2.Object) error) error {
	field, ok := u.Object["items"]
	if !ok {
		return errors.New("content is not a list")
	}
	items, ok := field.([]interface{})
	if !ok {
		return fmt.Errorf("content is not a list: %T", field)
	}
	for _, item := range items {
		child, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("items member is not an object: %T", child)
		}
		if err := fn(&Unstructured{Object: child}); err != nil {
			return err
		}
	}
	return nil
}

func (u *Unstructured) GetAPIVersion() string {
	return getNestedString(u.Object, "apiVersion")
}

func (u *Unstructured) SetAPIVersion(version string) {
	u.setNestedField(version, "apiVersion")
}

func (u *Unstructured) GetKind() string {
	return getNestedString(u.Object, "kind")
}

func (u *Unstructured) SetKind(kind string) {
	u.setNestedField(kind, "kind")
}

func (u *Unstructured) DeepCopy() *Unstructured {
	if u == nil {
		return nil
	}
	out := new(Unstructured)
	*out = *u
	out.Object = runtime.DeepCopyJSON(u.Object)
	return out
}

func (u *Unstructured) setNestedField(value interface{}, fields ...string) {
	if u.Object == nil {
		u.Object = make(map[string]interface{})
	}
	SetNestedField(u.Object, value, fields...)
}
