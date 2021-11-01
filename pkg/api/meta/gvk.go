package meta

import (
	"fmt"
	"strings"
)

// GroupKind specifies a Group and a Kind, but does not force a version.  This is useful for identifying
// concepts during lookup stages without having partially valid types_old
type GroupKind struct {
	Group string
	Kind  string
}

func (gk GroupKind) Empty() bool {
	return len(gk.Group) == 0 && len(gk.Kind) == 0
}

func (gk GroupKind) WithVersion(version string) GroupVersionKind {
	return GroupVersionKind{Group: gk.Group, Version: version, Kind: gk.Kind}
}

func (gk GroupKind) String() string {
	if len(gk.Group) == 0 {
		return gk.Kind
	}
	return gk.Kind + "." + gk.Group
}

// GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// Empty returns true if group, version, and kind are empty
func (gvk GroupVersionKind) Empty() bool {
	return len(gvk.Group) == 0 && len(gvk.Version) == 0 && len(gvk.Kind) == 0
}

func (gvk GroupVersionKind) GroupKind() GroupKind {
	return GroupKind{Group: gvk.Group, Kind: gvk.Kind}
}

func (gvk GroupVersionKind) GroupVersion() GroupVersion {
	return GroupVersion{Group: gvk.Group, Version: gvk.Version}
}

func (gvk GroupVersionKind) String() string {
	return gvk.Group + "/" + gvk.Version + ", Kind=" + gvk.Kind
}

// GroupVersion contains the "group" and the "version", which uniquely identifies the API.
type GroupVersion struct {
	Group   string
	Version string
}

// Empty returns true if group and version are empty
func (gv GroupVersion) Empty() bool {
	return len(gv.Group) == 0 && len(gv.Version) == 0
}

// String puts "group" and "version" into a single "group/version" string. For the legacy v1
// it returns "v1".
func (gv GroupVersion) String() string {
	if len(gv.Group) > 0 {
		return gv.Group + "/" + gv.Version
	}
	return gv.Version
}

// Identifier implements runtime.GroupVersioner interface.
func (gv GroupVersion) Identifier() string {
	return gv.String()
}

// KindForGroupVersionKinds identifies the preferred GroupVersionKind out of a list. It returns ok false
// if none of the options match the group. It prefers a match to group and version over just group.
// TODO: Move GroupVersion to a package under pkg/runtime, since it's used by scheme.
// TODO: Introduce an adapter type between GroupVersion and runtime.GroupVersioner, and use LegacyCodec(GroupVersion)
//   in fewer places.
func (gv GroupVersion) KindForGroupVersionKinds(kinds []GroupVersionKind) (target GroupVersionKind, ok bool) {
	for _, gvk := range kinds {
		if gvk.Group == gv.Group && gvk.Version == gv.Version {
			return gvk, true
		}
	}
	for _, gvk := range kinds {
		if gvk.Group == gv.Group {
			return gv.WithKind(gvk.Kind), true
		}
	}
	return GroupVersionKind{}, false
}

// ParseGroupVersion turns "group/version" string into a GroupVersion struct. It reports error
// if it cannot parse the string.
func ParseGroupVersion(gv string) (GroupVersion, error) {
	// this can be the internal version for the legacy kube types_old
	// TODO once we've cleared the last uses as strings, this special case should be removed.
	if (len(gv) == 0) || (gv == "/") {
		return GroupVersion{}, nil
	}

	switch strings.Count(gv, "/") {
	case 0:
		return GroupVersion{"", gv}, nil
	case 1:
		i := strings.Index(gv, "/")
		return GroupVersion{gv[:i], gv[i+1:]}, nil
	default:
		return GroupVersion{}, fmt.Errorf("unexpected GroupVersion string: %v", gv)
	}
}

// WithKind creates a GroupVersionKind based on the method receiver's GroupVersion and the passed Kind.
func (gv GroupVersion) WithKind(kind string) GroupVersionKind {
	return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
}

// WithResource creates a GroupVersionResource based on the method receiver's GroupVersion and the passed Resource.
func (gv GroupVersion) WithResource(resource string) GroupVersionResource {
	return GroupVersionResource{Group: gv.Group, Version: gv.Version, Resource: resource}
}

// GroupVersions can be used to represent a set of desired group versions.
// TODO: Move GroupVersions to a package under pkg/runtime, since it's used by scheme.
// TODO: Introduce an adapter type between GroupVersions and runtime.GroupVersioner, and use LegacyCodec(GroupVersion)
//   in fewer places.
type GroupVersions []GroupVersion

// Identifier implements runtime.GroupVersioner interface.
func (gvs GroupVersions) Identifier() string {
	groupVersions := make([]string, 0, len(gvs))
	for i := range gvs {
		groupVersions = append(groupVersions, gvs[i].String())
	}
	return fmt.Sprintf("[%s]", strings.Join(groupVersions, ","))
}

// KindForGroupVersionKinds identifies the preferred GroupVersionKind out of a list. It returns ok false
// if none of the options match the group.
func (gvs GroupVersions) KindForGroupVersionKinds(kinds []GroupVersionKind) (GroupVersionKind, bool) {
	var targets []GroupVersionKind
	for _, gv := range gvs {
		target, ok := gv.KindForGroupVersionKinds(kinds)
		if !ok {
			continue
		}
		targets = append(targets, target)
	}
	if len(targets) == 1 {
		return targets[0], true
	}
	if len(targets) > 1 {
		return bestMatch(kinds, targets), true
	}
	return GroupVersionKind{}, false
}

// bestMatch tries to pick best matching GroupVersionKind and falls back to the first
// found if no exact match exists.
func bestMatch(kinds []GroupVersionKind, targets []GroupVersionKind) GroupVersionKind {
	for _, gvk := range targets {
		for _, k := range kinds {
			if k == gvk {
				return k
			}
		}
	}
	return targets[0]
}

// ToAPIVersionAndKind is a convenience method for satisfying runtime.Object on types_old that
// do not use TypeMeta.
func (gvk GroupVersionKind) ToAPIVersionAndKind() (string, string) {
	if gvk.Empty() {
		return "", ""
	}
	return gvk.GroupVersion().String(), gvk.Kind
}

// FromAPIVersionAndKind returns a GVK representing the provided fields for types_old that
// do not use TypeMeta. This method exists to support test types_old and legacy serializations
// that have a distinct group and kind.
// TODO: further reduce usage of this method.
func FromAPIVersionAndKind(apiVersion, kind string) GroupVersionKind {
	if gv, err := ParseGroupVersion(apiVersion); err == nil {
		return GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}
	}
	return GroupVersionKind{Kind: kind}
}

type GroupResourcer interface {
	GroupResource() GroupResource
}

type GroupResource struct {
	Group    string
	Resource string
}

func (gr GroupResource) WithVersion(version string) GroupVersionResource {
	return GroupVersionResource{Group: gr.Group, Version: version, Resource: gr.Resource}
}

func (gr GroupResource) String() string {
	if len(gr.Group) == 0 {
		return gr.Resource
	}
	return gr.Resource + "." + gr.Group
}

func (gr GroupResource) GroupResource() GroupResource {
	return gr
}

type GroupVersionResource struct {
	Group    string
	Resource string
	Version  string
}

func (g GroupVersionResource) String() string {
	return fmt.Sprintf("%s/%s, Resource=%s", g.Group, g.Version, g.Resource)
}

func (g GroupVersionResource) GroupResource() GroupResource {
	return GroupResource{Group: g.Group, Resource: g.Resource}
}

func (g GroupVersionResource) GroupVersion() GroupVersion {
	return GroupVersion{Group: g.Group, Version: g.Version}
}
