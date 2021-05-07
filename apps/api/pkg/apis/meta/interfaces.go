package meta

import (
	"fmt"
	metav1 "github.com/nrc-no/core/apps/api/pkg/apis/meta/v1"
	"github.com/nrc-no/core/apps/api/pkg/runtime"
)

type MetadataAccessor interface {
	APIVersion(obj runtime.Object) (string, error)
	SetAPIVersion(obj runtime.Object, version string) error

	Kind(obj runtime.Object) (string, error)
	SetKind(obj runtime.Object, kind string) error

	UID(obj runtime.Object) (string, error)
	SetUID(obj runtime.Object, uid string) error

	Labels(obj runtime.Object) (map[string]string, error)
	SetLabels(obj runtime.Object, labels map[string]string) error

	Annotations(obj runtime.Object) (map[string]string, error)
	SetAnnotations(obj runtime.Object, annotations map[string]string) error

	ResourceVersioner
}

func NewAccessor() MetadataAccessor {
	return resourceAccessor{}
}

type ResourceVersioner interface {
	ResourceVersion(obj runtime.Object) (int, error)
	SetResourceVersion(obj runtime.Object, version int) error
}

type resourceAccessor struct{}

var _ MetadataAccessor = &resourceAccessor{}

func (r resourceAccessor) APIVersion(obj runtime.Object) (string, error) {
	return objectAccessor{obj}.GetAPIVersion(), nil
}

func (r resourceAccessor) SetAPIVersion(obj runtime.Object, version string) error {
	objectAccessor{obj}.SetAPIVersion(version)
	return nil
}

func (r resourceAccessor) Kind(obj runtime.Object) (string, error) {
	return objectAccessor{obj}.GetKind(), nil
}

func (r resourceAccessor) SetKind(obj runtime.Object, kind string) error {
	objectAccessor{obj}.SetKind(kind)
	return nil
}

func (r resourceAccessor) UID(obj runtime.Object) (string, error) {
	accessor, err := Accessor(obj)
	if err != nil {
		return "", err
	}
	return accessor.GetUID(), nil
}

func (r resourceAccessor) SetUID(obj runtime.Object, uid string) error {
	accessor, err := Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetUID(uid)
	return nil
}

func (r resourceAccessor) Labels(obj runtime.Object) (map[string]string, error) {
	accessor, err := Accessor(obj)
	if err != nil {
		return nil, err
	}
	return accessor.GetLabels(), nil
}

func (r resourceAccessor) SetLabels(obj runtime.Object, labels map[string]string) error {
	accessor, err := Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetLabels(labels)
	return nil
}

func (r resourceAccessor) Annotations(obj runtime.Object) (map[string]string, error) {
	accessor, err := Accessor(obj)
	if err != nil {
		return nil, err
	}
	return accessor.GetAnnotations(), nil
}

func (r resourceAccessor) SetAnnotations(obj runtime.Object, annotations map[string]string) error {
	accessor, err := Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetAnnotations(annotations)
	return nil
}

func (r resourceAccessor) ResourceVersion(obj runtime.Object) (int, error) {
	accessor, err := Accessor(obj)
	if err != nil {
		return 0, err
	}
	return accessor.GetResourceVersion(), nil
}

func (r resourceAccessor) SetResourceVersion(obj runtime.Object, version int) error {
	accessor, err := Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetResourceVersion(version)
	return nil
}

type objectAccessor struct {
	runtime.Object
}

func (obj objectAccessor) GetKind() string {
	return obj.GetKind()
}

func (obj objectAccessor) SetKind(kind string) {
	obj.SetKind(kind)
}

func (obj objectAccessor) GetAPIVersion() string {
	return obj.GetAPIVersion()
}

func (obj objectAccessor) SetAPIVersion(version string) {
	obj.SetAPIVersion(version)
}

var errNotObject = fmt.Errorf("object does not implement the Object interfaces")

func Accessor(obj interface{}) (metav1.Object, error) {
	switch t := obj.(type) {
	case metav1.Object:
		return t, nil
	case metav1.ObjectMetaAccessor:
		return t.GetObjectMeta(), nil
	default:
		return nil, errNotObject
	}
}
