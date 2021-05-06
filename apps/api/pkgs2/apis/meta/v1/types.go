package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ObjectMeta struct {
	Name              string            `json:"name,omitempty"`
	UID               string            `json:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
	DeletionTimestamp time.Time         `json:"deletionTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations"`
}

type APIResource struct {
	Name         string `json:"name"`
	SingularName string `json:"singularName"`
	Group        string `json:"group"`
	Version      string `json:"version"`
	Kind         string `json:"kind"`
}

type CreateOptions struct {
	TypeMeta `json:",inline"`
}

type ListOptions struct {
	metav1.TypeMeta
}
