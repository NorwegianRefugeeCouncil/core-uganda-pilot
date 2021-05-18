package storageversion

import (
	"github.com/nrc-no/core/apps/api/pkg/runtime"
	"github.com/nrc-no/core/apps/api/pkg/runtime/schema"
)

// ResourceInfo contains the information to register the resource to the
// storage version API.
type ResourceInfo struct {
	GroupResource schema.GroupResource

	EncodingVersion string
	// Used to calculate decodable versions. Can only be used after all
	// equivalent versions are registered by InstallREST.
	EquivalentResourceMapper runtime.EquivalentResourceRegistry

	// DirectlyDecodableVersions is a list of versions that the converter for REST storage knows how to convert.  This
	// contains items like apiextensions.k8s.io/v1beta1 even if we don't serve that version.
	DirectlyDecodableVersions []schema.GroupVersion
}
