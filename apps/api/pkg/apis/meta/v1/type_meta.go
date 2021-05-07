package v1

import "strings"

// TypeMeta represents an individual object in an API response or request.
// It represents the API schema version and kind/type of object
type TypeMeta struct {

	// APIVersion defines the versioned schema of this representation
	// of an object.
	APIVersion string `json:"apiVersion,omitempty" bson:"apiVersion,omitempty"`

	// Kind is a string value representing the REST resource this
	// object represents
	Kind string `json:"kind,omitempty" bson:"kind,omitempty"`
}

func (t *TypeMeta) GetAPIVersion() string {
	return t.APIVersion
}

func (t *TypeMeta) SetAPIVersion(version string) {
	t.APIVersion = version
}

func (t *TypeMeta) GetAPIGroup() string {
	parts := strings.Split(t.APIVersion, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func (t *TypeMeta) SetAPIGroup(group string) {
	if len(t.APIVersion) == 0 {
		t.APIVersion = group + "/"
		return
	}
	parts := strings.Split(t.APIVersion, "/")
	t.APIVersion = group + "/" + parts[1]
}

func (t *TypeMeta) GetKind() string {
	return t.Kind
}

func (t *TypeMeta) SetKind(kind string) {
	t.Kind = kind
}
