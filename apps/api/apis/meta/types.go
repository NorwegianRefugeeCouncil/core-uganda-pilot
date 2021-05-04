package meta

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//go:generate controller-gen object paths=$GOFILE

type APIResource struct {
  Name         string `json:"name"`
  SingularName string `json:"singularName"`
  Group        string `json:"group"`
  Version      string `json:"version"`
  Kind         string `json:"kind"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ListOptions struct{
  metav1.TypeMeta
}

