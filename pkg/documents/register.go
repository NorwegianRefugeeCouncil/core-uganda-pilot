package documents

import "github.com/nrc-no/core/pkg/api/meta"

const GroupName = "doc.nrc.no"
const Version = "v1"

var GroupVersion = meta.GroupVersion{Group: GroupName, Version: Version}

func Kind(kind string) meta.GroupKind {
	return GroupVersion.WithKind(kind).GroupKind()
}
