package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type SimpleObjectMeta struct {
	Name              string            `json:"name"`
	ResourceVersion   string            `json:"resourceVersion"`
	UID               types.UID         `json:"uid"`
	Generation        int64             `json:"generation"`
	CreationTimestamp metav1.Time       `json:"creationTimestamp"`
	DeletionTimestamp *metav1.Time      `json:"deletionTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

func (s *SimpleObjectMeta) GetResourceVersion() string {
	return s.ResourceVersion
}

func (s *SimpleObjectMeta) SetResourceVersion(version string) {
	s.ResourceVersion = version
}

func (s *SimpleObjectMeta) GetSelfLink() string {
	return ""
}

func (s *SimpleObjectMeta) SetSelfLink(link string) {
	return
}

func (s *SimpleObjectMeta) GetNamespace() string {
	return ""
}

func (s *SimpleObjectMeta) SetNamespace(namespace string) {
	return
}

func (s *SimpleObjectMeta) GetName() string {
	return s.Name
}

func (s *SimpleObjectMeta) SetName(name string) {
	s.Name = name
}

func (s *SimpleObjectMeta) GetGenerateName() string {
	return ""
}

func (s *SimpleObjectMeta) SetGenerateName(name string) {

}

func (s *SimpleObjectMeta) GetUID() types.UID {
	return s.UID
}

func (s *SimpleObjectMeta) SetUID(uid types.UID) {
	s.UID = uid
}

func (s *SimpleObjectMeta) GetGeneration() int64 {
	return s.Generation
}

func (s *SimpleObjectMeta) SetGeneration(generation int64) {
	s.Generation = generation
}

func (s *SimpleObjectMeta) GetCreationTimestamp() metav1.Time {
	return s.CreationTimestamp
}

func (s *SimpleObjectMeta) SetCreationTimestamp(timestamp metav1.Time) {
	s.CreationTimestamp = timestamp
}

func (s *SimpleObjectMeta) GetDeletionTimestamp() *metav1.Time {
	return s.DeletionTimestamp
}

func (s *SimpleObjectMeta) SetDeletionTimestamp(timestamp *metav1.Time) {
	s.DeletionTimestamp = timestamp
}

func (s *SimpleObjectMeta) GetDeletionGracePeriodSeconds() *int64 {
	return nil
}

func (s *SimpleObjectMeta) SetDeletionGracePeriodSeconds(*int64) {
	return
}

func (s *SimpleObjectMeta) GetLabels() map[string]string {
	return s.Labels
}

func (s *SimpleObjectMeta) SetLabels(labels map[string]string) {
	s.Labels = labels
}

func (s *SimpleObjectMeta) GetAnnotations() map[string]string {
	return s.Annotations
}

func (s *SimpleObjectMeta) SetAnnotations(annotations map[string]string) {
	s.Annotations = annotations
}

func (s *SimpleObjectMeta) GetFinalizers() []string {
	return []string{}
}

func (s *SimpleObjectMeta) SetFinalizers(finalizers []string) {
	return
}

func (s *SimpleObjectMeta) GetOwnerReferences() []metav1.OwnerReference {
	return []metav1.OwnerReference{}
}

func (s *SimpleObjectMeta) SetOwnerReferences([]metav1.OwnerReference) {
	return
}

func (s *SimpleObjectMeta) GetClusterName() string {
	return ""
}

func (s *SimpleObjectMeta) SetClusterName(clusterName string) {
	return
}

func (s *SimpleObjectMeta) GetManagedFields() []metav1.ManagedFieldsEntry {
	return []metav1.ManagedFieldsEntry{}
}

func (s *SimpleObjectMeta) SetManagedFields(managedFields []metav1.ManagedFieldsEntry) {
	return
}

func (s *SimpleObjectMeta) GetObjectMeta() metav1.Object {
	return s
}
