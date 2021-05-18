package v1

import (
  "github.com/nrc-no/core/apps/api/pkg/runtime"
  "github.com/nrc-no/core/apps/api/pkg/runtime/schema"
  utilruntime "github.com/nrc-no/core/apps/api/pkg/util/runtime"
)

const GroupName = "meta"

var (
  schemeBuilder      runtime.SchemeBuilder
  localSchemeBuilder = &schemeBuilder
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

var Unversioned = schema.GroupVersion{Group: "", Version: "v1"}

// WatchEventKind is name reserved for serializing watch events.
const WatchEventKind = "WatchEvent"

// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
  return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// scheme is the registry for the common types that adhere to the meta v1 API spec.
var scheme = runtime.NewScheme()

// ParameterCodec knows about query parameters used with the meta v1 API spec.
var ParameterCodec = runtime.NewParameterCodec(scheme)

var optionsTypes = []runtime.Object{
  &ListOptions{},
  &GetOptions{},
  &DeleteOptions{},
  &CreateOptions{},
  &UpdateOptions{},
}

func AddToGroupVersion(scheme *runtime.Scheme, groupVersion schema.GroupVersion) {
  scheme.AddKnownTypeWithName(groupVersion.WithKind(WatchEventKind), &WatchEvent{})
  scheme.AddKnownTypeWithName(
    schema.GroupVersion{Group: groupVersion.Group, Version: runtime.APIVersionInternal}.WithKind(WatchEventKind),
    &InternalEvent{},
  )
  scheme.AddKnownTypes(groupVersion, optionsTypes...)
  scheme.AddUnversionedTypes(Unversioned,
    &Status{},
    &APIVersions{},
    &APIGroupList{},
    &APIGroup{},
    &APIResourceList{},
  )
  utilruntime.Must(RegisterConversions(scheme))
}

// Unlike other API groups, meta internal knows about all meta external versions, but keeps
// the logic for conversion private.
func init() {
  scheme.AddKnownTypes(SchemeGroupVersion, optionsTypes...)
}
